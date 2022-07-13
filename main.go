package main

import (
	"crypto/ed25519"
	"encoding/base64"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"
)

type AppConfig struct {
	Port       string `json:"Port"`
	KeySet     string `json:"KeySet"`
	PrivateKey string `json:"PrivateKey"`
}

type HttpResponse struct {
	SignedURL string `json:"SignedURL"`
}

func main() {
	fmt.Println("Starting the Application..")

	//Initialize App Configuration
	config := AppConfig{
		Port: "8080",
	}
	config.KeySet = os.Getenv("KEYSET")
	config.PrivateKey = os.Getenv("PRIVATEKEY")

	// HTTP Handlers
	http.HandleFunc("/signurl", http_signurl(config))

	// Start http server
	fmt.Printf("Starting Listening on port %v..\n", config.Port)
	http.ListenAndServe(":"+config.Port, nil)

}

func http_signurl(config AppConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		url := r.URL.Query().Get("url")
		expiration := r.URL.Query().Get("expiration")

		if url == "" || expiration == "" {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "Bad Request, please provide \"url\" parameter")
			return
		}

		tExpiration, err := stringToUnix(expiration)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "Error: %s", err)
			return
		}

		resp := HttpResponse{}

		SignedURL, err := signUrl(config, url, tExpiration)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "Error: %s", err)
			return
		}

		resp.SignedURL = SignedURL

		fmt.Fprintf(w, "%+v", resp)
	}
}

func stringToUnix(s string) (t time.Time, e error) {
	i, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return t, fmt.Errorf("Couldn't parse expiration parameter: %s as Unix Time", s)
	}

	tm := time.Unix(i, 0)

	return tm, nil
}

func signUrl(config AppConfig, url string, expiration time.Time) (signedurl string, err error) {
	keyset, err := base64.RawStdEncoding.DecodeString(config.KeySet)
	if err != nil {
		return signedurl, fmt.Errorf("Could not parse the private key: %s", err)
	}

	toSign := fmt.Sprintf("%s?Expires=%d&KeyName=%s", url, expiration.Unix(), config.KeySet)
	signedurl = string(ed25519.Sign(keyset, []byte(toSign)))

	return signedurl, nil

}
