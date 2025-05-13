// Package keystore implements the auth.KeyLookup interface. This implements
// an in-memory keystore for JWT support.
// The whole point is to implement the interface for retrieving keys.
package keystore

import (
	"bytes"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"path"
	"strings"

	"github.com/golang-jwt/jwt/v4"
)

// PrivateKey represents key information.
// Holds both the parsed *rsa.PrivateKey and the original PEM bytes.
type PrivateKey struct {
	PK  *rsa.PrivateKey
	PEM []byte
}

// KeyStore represents an in memory store implementation of the
// KeyLookup interface for use with the auth package.
type KeyStore struct {
	store map[string]PrivateKey
}

// New constructs an empty KeyStore ready for use.
func New() *KeyStore {
	return &KeyStore{
		store: make(map[string]PrivateKey),
	}
}

// NewMap constructs a KeyStore with an initial set of keys.
func NewMap(store map[string]PrivateKey) *KeyStore {
	return &KeyStore{
		store: store,
	}
}

// NewFS constructs a KeyStore based on a set of PEM files rooted inside
// of a directory. The name of each PEM file will be used as the key id.
// Example: keystore.NewFS(os.DirFS("/zarf/keys/"))
// Example: /zarf/keys/54bb2165-71e1-41a6-af3e-7da4a0e1e2c1.pem
// fsys fs.FS is any file‐system abstraction (e.g., os.DirFS("/keys") or an embedded FS).
func NewFS(fsys fs.FS) (*KeyStore, error) {
	ks := New()
	// Define a walk‐callback fn that will be called for every file or directory under fsys.
	// It receives: fileName (path relative to the root, e.g. "54bb…c1.pem")
	// dirEntry (info on whether it’s a file or directory) & err from trying to read that directory entry.
	fn := func(fileName string, dirEntry fs.DirEntry, err error) error {
		if err != nil {
			return fmt.Errorf("walkdir failure: %w", err)
		}
		// Skip directories – we only care about files.
		if dirEntry.IsDir() {
			return nil
		}
		// Skip non-PEM files – only process files ending in .pem.
		if path.Ext(fileName) != ".pem" {
			return nil
		}

		file, err := fsys.Open(fileName)
		if err != nil {
			return fmt.Errorf("opening key file: %w", err)
		}
		defer file.Close()

		// limit PEM file size to 1 megabyte. This should be reasonable for
		// almost any PEM file and prevents shenanigans like linking the file
		// to /dev/random or something like that.
		pem, err := io.ReadAll(io.LimitReader(file, 1024*1024))
		if err != nil {
			return fmt.Errorf("reading auth private key: %w", err)
		}

		pk, err := jwt.ParseRSAPrivateKeyFromPEM(pem)
		if err != nil {
			return fmt.Errorf("parsing auth private key: %w", err)
		}

		key := PrivateKey{
			PK:  pk,
			PEM: pem,
		}

		ks.store[strings.TrimSuffix(dirEntry.Name(), ".pem")] = key

		return nil
	}
	// Invoke fs.WalkDir, starting at "." (the root of fsys), calling our fn for every entry.
	if err := fs.WalkDir(fsys, ".", fn); err != nil {
		return nil, fmt.Errorf("walking directory: %w", err)
	}

	return ks, nil
}

// PrivateKey searches the key store for a given kid and returns the private key.
func (ks *KeyStore) PrivateKey(kid string) (string, error) {
	privateKey, found := ks.store[kid]
	if !found {
		return "", errors.New("kid lookup failed")
	}

	return string(privateKey.PEM), nil
}

// PublicKey searches the key store for a given kid and returns the public key.
func (ks *KeyStore) PublicKey(kid string) (string, error) {
	privateKey, found := ks.store[kid]
	if !found {
		return "", errors.New("kid lookup failed")
	}

	asn1Bytes, err := x509.MarshalPKIXPublicKey(&privateKey.PK.PublicKey)
	if err != nil {
		return "", fmt.Errorf("marshaling public key: %w", err)
	}

	block := pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: asn1Bytes,
	}

	var b bytes.Buffer
	if err := pem.Encode(&b, &block); err != nil {
		return "", fmt.Errorf("encoding to private file: %w", err)
	}

	return b.String(), nil
}
