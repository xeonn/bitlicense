package issuer

import (
	"crypto/ed25519"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"

	"golang.org/x/crypto/ssh"
)

func ValidateFile(filePath string) bool {
	license, err := loadFile(filePath)
	if err != nil {
		fmt.Println(err)
		return false
	}

	return Validate(license)
}

func Validate(license *License) bool {
	return validate(license, "id_ed25519.pub")
}

func validate(license *License, keyPath string) bool {
	// load public key
	// hashContent client + expiry date
	// decrypt signature
	// compare decrypted signature with hashContent

	publicKey, err := unmarshalPublicKey("certs/" + keyPath)
	if err != nil {
		fmt.Println(err)
		return false
	}

	byteHash := hashContent(license.Client, license.Expiry)

	// calculatedHash := base64.StdEncoding.EncodeToString(byteHash)
	decryptedHash, err := base64.StdEncoding.DecodeString(license.Signature)
	if err != nil {
		fmt.Println(err)
		return false
	}

	// decrypt expiry date
	return ed25519.Verify(*publicKey, byteHash, decryptedHash)
}


// read and return public key from file
func unmarshalPublicKey(keyPath string) (*ed25519.PublicKey, error) {
	pemBytes, err := os.ReadFile(keyPath)
	if err != nil {
		return nil, err
	}

	parsed, _, _, _, err := ssh.ParseAuthorizedKey(pemBytes)
	if err != nil {
		return nil, err
	}

	// upgrade to ssh.CryptoPublicKey interface
	sshkey := parsed.(ssh.CryptoPublicKey)
	edKey := sshkey.CryptoPublicKey().(ed25519.PublicKey)

	return &edKey, nil
}

func loadFile(filePath string) (*License, error) {
	// read file
	// parse json
	// return License struct

	// read file
	file, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	// parse json
	var license License
	err = json.Unmarshal(file, &license)
	if err != nil {
		return nil, err
	}

	return &license, nil
}

