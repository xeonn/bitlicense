package issuer

import (
	"crypto/ed25519"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"os"
	"time"

	"golang.org/x/crypto/ssh"
)

// Issue license file
func Issue(client string, expiryDate time.Time) *License {
	// Mon Jan 2 15:04:05 MST 2006
	return issue(client, expiryDate.UTC().Format(time.RFC3339), "privatekey")
}

func issue(client string, expiryDate string, keyPath string) *License {
	// load private key
	// hash client + expiry date
	// encrypt/sign expiry date
	// return License struct

	privateKey, err := unmarshalPrivateKey("certs/" + keyPath)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	// signature consists of Client:ExpiryDate string concatenated with semicolon
	hashedContent := hashContent(client, expiryDate)

	// sign(encrypt) the hashed content and convert to base64
	signature := base64.StdEncoding.EncodeToString(ed25519.Sign(*privateKey, hashedContent))

	return &License{
		Client:    client,
		Expiry:    expiryDate,
		Signature: signature,
	}
}

func hashContent(client string, expirydate string) []byte {
	hasher := sha256.New()
	hasher.Write([]byte(client + ":" + expirydate))
	return hasher.Sum(nil)
}

// read and return private key from file
func unmarshalPrivateKey(keyPath string) (*ed25519.PrivateKey, error) {
	pemBytes, err := os.ReadFile(keyPath)
	if err != nil {
		return nil, err
	}

	privateKey, err := ssh.ParseRawPrivateKey(pemBytes)
	if err != nil {
		return nil, err
	}

	return privateKey.(*ed25519.PrivateKey), nil
}
