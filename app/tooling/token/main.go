package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	_ "embed"
	"encoding/pem"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

func main() {
	if err := genToken(); err != nil {
		log.Fatalln(err)
	}
}

func genToken() error {
	privateKey, err := getPrivateKey()
	if err != nil {
		return fmt.Errorf("failed to get Private Key %w", err)
	}
	claims := struct {
		jwt.RegisteredClaims
		Roles []string
	}{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   "1234",
			Issuer:    "URL-Shortener",
			ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(8760 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
		},
		Roles: []string{"ADMIN"},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	token.Header["kid"] = "54bb2165-71e1-41a6-af3e-7da4a0e1e2c1"

	str, err := token.SignedString(privateKey)
	if err != nil {
		return fmt.Errorf("signing token: %w", err)
	}
	fmt.Println("*******************")
	fmt.Println(str)
	fmt.Println("*******************")
	// =========================================================================
	// Validate token
	// Only tokens signed with RS256 will be considered
	parser := jwt.NewParser(jwt.WithValidMethods([]string{jwt.SigningMethodRS256.Name}))

	var claims2 struct {
		jwt.RegisteredClaims
		Roles []string
	}
	// it takes the token as the first argument, the claims that we want to pull back out of the token.
	// also it doesn't take the public key directly it takes a key function that knows how to look up what public key
	// we should be using.
	// In a more robust setup you'd inspect "t.Header["kid"]" to choose among multiple keys
	keyFunc := func(t *jwt.Token) (interface{}, error) {
		return &privateKey.PublicKey, nil
	}

	tkn, err := parser.ParseWithClaims(str, &claims2, keyFunc)

	if err != nil {
		return fmt.Errorf("parsing token: %w", err)
	}

	if !tkn.Valid {
		return fmt.Errorf("token not valid")
	}

	fmt.Println("SIGNATURE VALIDATED")
	fmt.Printf("%#v", claims2)
	fmt.Printf("Roles %#v", claims2.Roles)
	fmt.Println("*******************")
	return nil
}

func getPrivateKey() (*rsa.PrivateKey, error) {
	keysDirPath := filepath.Join("zarf", "keys")
	// 1) Stat the directory
	if _, err := os.Stat(keysDirPath); err != nil {
		if os.IsNotExist(err) {
			// 2a) It doesn't exist → create it (and any parent paths)
			if mkErr := os.MkdirAll(keysDirPath, 0o755); mkErr != nil {
				return nil, fmt.Errorf("failed to create keys directory %q: %w",
					keysDirPath, mkErr)
			}
		} else {
			// 2b) Some other error (e.g. permissions) → bail out
			return nil, fmt.Errorf("could not stat keys directory %q: %w",
				keysDirPath, err)
		}
	}
	// Path is relative to the project root, since I am running the main from the root folder
	privatekeyFilePath := filepath.Join("zarf", "keys", "54bb2165-71e1-41a6-af3e-7da4a0e1e2c1.pem")

	// Check if the file exists by opening it
	privateKeyFile, err := os.Open(privatekeyFilePath)
	if err != nil {
		fmt.Println("Error trying to open private key file, err:", err)
		fmt.Println("Generating private & public keys ...")
		if err := genKeyFiles(); err != nil {
			return nil, err
		}
		privateKeyFile, err = os.Open(privatekeyFilePath)
		if err != nil {
			return nil, err
		}
	}
	defer privateKeyFile.Close()

	privateKeyBytes, err := io.ReadAll(privateKeyFile)
	if err != nil {
		fmt.Println("Error trying to read private key bytes, err:", err)
		return nil, err
	}

	block, _ := pem.Decode(privateKeyBytes)
	if block == nil {
		return nil, fmt.Errorf("no PEM block found in %q", privatekeyFilePath)
	}

	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse private key bytes %w", err)
	}
	return privateKey, nil
}

func genKeyFiles() error {
	// Generate a new private key.
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return fmt.Errorf("generating key: %w", err)
	}

	// Create a file for the private key information in PEM form.
	privateKeyFilePath := filepath.Join("zarf", "keys", "54bb2165-71e1-41a6-af3e-7da4a0e1e2c1.pem")
	privateFile, err := os.Create(privateKeyFilePath)
	if err != nil {
		return fmt.Errorf("creating private file: %w", err)
	}
	defer privateFile.Close()
	// Here we are encoding it into the pem formatting to write to disk
	// Construct a PEM block for the private key.
	privateBlock := pem.Block{
		Type: "PRIVATE KEY",
		// We have to use this "x509.MarshalPKCS1PrivateKey" function to do this properly.
		Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
	}

	// Write the private key to the private key file.
	if err := pem.Encode(privateFile, &privateBlock); err != nil {
		return fmt.Errorf("encoding to private file: %w", err)
	}

	fmt.Println("private key files generated")

	return nil
}
