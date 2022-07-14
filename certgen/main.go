package main

import (
	"crypto/ed25519"
	"encoding/base64"
	"fmt"
	"os"
)

func main() {
	keyFile, err := os.OpenFile("private.key", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		fmt.Fprintf(os.Stderr, "could not open private key file for writing: %s\n", err)
		os.Exit(1)
	}
	pubFile, err := os.OpenFile("public.key", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		fmt.Fprintf(os.Stderr, "could not open public key file for writing: %s\n", err)
		os.Exit(1)
	}

	pubKey, privKey, err := ed25519.GenerateKey( /*rand=*/ nil)
	if err != nil {
		fmt.Printf("could not generate key: %v\n", err)
		os.Exit(1)
	}

	if _, err := keyFile.WriteString(base64.RawURLEncoding.EncodeToString(privKey)); err != nil {
		fmt.Printf("could not write private key: %v\n", err)
		os.Exit(1)
	}

	if _, err := pubFile.WriteString(base64.RawURLEncoding.EncodeToString(pubKey)); err != nil {
		fmt.Printf("could not write public key: %v\n", err)
		os.Exit(1)
	}

}
