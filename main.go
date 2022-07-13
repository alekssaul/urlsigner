package main

import (
	"crypto/ed25519"
	"encoding/base64"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

type AppConfig struct {
	Port          string `json:"Port"`
	KeySet        string `json:"KeySet"`
	PrivateKey    string `json:"PrivateKey"`
	binPrivateKey []byte
}

type HttpResponse struct {
	SignedURL string `json:"SignedURL"`
}

func main() {
	fmt.Println("Starting the Application..")

	config := initApp()

	// HTTP Handlers
	http.HandleFunc("/signurl", http_signurl(config))

	// Start http server
	fmt.Printf("Starting Listening on port %v..\n", config.Port)
	http.ListenAndServe(":"+config.Port, nil)

}

func initApp() (config AppConfig) {
	//Initialize App Configuration
	config.Port = os.Getenv("PORT")
	if config.Port == "" {
		config.Port = "8080"
	}
	config.KeySet = os.Getenv("KEYSET")
	if config.KeySet == "" {
		fmt.Printf("Must set \"KEYSET\" Environmental Variable\n")
		os.Exit(1)
	}

	config.PrivateKey = os.Getenv("PRIVATEKEY")
	if config.PrivateKey == "" {
		fmt.Printf("Must set \"PRIVATEKEY\" Environmental Variable\n")
		os.Exit(1)
	}

	keyset, err := base64.RawURLEncoding.DecodeString(config.PrivateKey)
	if err != nil {
		fmt.Printf("Could not parse the private key: %s\n", err)
		os.Exit(1)
	}

	//keyset, _ := os.ReadFile("private.key")

	if len(keyset) != ed25519.PrivateKeySize {
		fmt.Printf("Private key Size is not valid. Expected: %v, Found %v\n", ed25519.PrivateKeySize, len(keyset))
		os.Exit(1)
	}

	config.binPrivateKey = keyset

	return config

}

func http_signurl(config AppConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		url := r.URL.Query().Get("url")
		expiration := r.URL.Query().Get("expiration")

		if url == "" || expiration == "" {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "Bad Request, please provide \"url\" and \"expiration\" parameters")
			return
		}

		tExpiration, err := stringToUnix(expiration)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "Error: %s", err)
			return
		}

		resp := HttpResponse{}

		SignedURL := signUrl(config, url, tExpiration)
		resp.SignedURL = SignedURL

		fmt.Fprintf(w, "%+v\n", resp)
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

func signUrl(config AppConfig, url string, expiration time.Time) (signedurl string) {
	sep := '?'
	if strings.ContainsRune(url, '?') {
		sep = '&'
	}
	toSign := fmt.Sprintf("%s%cExpires=%d&KeyName=%s", url, sep, expiration.Unix(), config.KeySet)
	sig := ed25519.Sign(config.binPrivateKey, []byte(toSign))
	return fmt.Sprintf("%s&Signature=%s", toSign, base64.RawURLEncoding.EncodeToString(sig))

}
