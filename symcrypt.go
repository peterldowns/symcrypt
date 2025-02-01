// symcrypt ("symmetric cryptography") provides a safe interface for symmetrically
// encrypting and decrypting values.
//
// Based on the examples at:
// - https://pkg.go.dev/golang.org/x/crypto/chacha20poly1305#example-NewX
// - https://cs.opensource.google/go/x/crypto/+/refs/tags/v0.24.0:chacha20poly1305/chacha20poly1305_test.go
//
// Providing `ownerID` is required for the "Authenticated" part of AEAD
// to actually function correctly. The idea is that if you encrypt a secret for
// one user, you want to prevent against a case where you accidentally decrypt
// it for a different user.  To do this, you send in the user's ID when you
// encrypt and decrypt messages.
//
// For more information, read https://en.wikipedia.org/wiki/Authenticated_encryption
package symcrypt

import (
	"crypto/cipher"
	cryptorand "crypto/rand"
	"encoding/hex"
	"io"

	// https://github.com/golang/crypto/blob/master/chacha20poly1305/chacha20poly1305.go
	"golang.org/x/crypto/chacha20poly1305"
)

// These types provide "Safety Through Incompatibility" —
// the goal is to make it harder to accidentally treat plaintext like
// ciphertext, or mix-and-match owners incorrectly.
//
// For more information on the general concept, you may enjoy reading:
// https://lukasschwab.me/blog/gen/safe-incompatibility.html
type (
	Plaintext  string
	Ciphertext string
	Owner      string
)

// Client allows for encrypting and decrypting secrets that are "owned" by
// something or someone that has a unique identifier. If you Encrypt() a
// plaintext value for a given owner, you will ownly be able to Decrypt() the
// ciphertext back into plaintext for that same owner. The goal of this
// interface is to make it harder for application programmers to accidentally
// decrypt secrets for someone other than the owner.
//
// For example usages, see the tests.
type Client interface {
	// Encrypt a Plaintext secret for a given Owner.
	Encrypt(plaintext Plaintext, ownerID Owner) (Ciphertext, error)
	// Decrypt a Ciphertext secret for a specific Owner.
	//
	// If the secret was not encrypted for this same Owner, the implementation
	// must return an error.
	Decrypt(ciphertext Ciphertext, ownerID Owner) (Plaintext, error)
}

var _ Client = &chachaClient{}

type chachaClient struct {
	aead cipher.AEAD
}

// A 32-byte (chacha20poly1305.KeySize) secret key, hex-encoded as a string.
type HexKey string

func NewClient(hexKey HexKey) (Client, error) {
	key, err := hex.DecodeString(string(hexKey))
	if err != nil {
		return nil, err
	}
	aead, err := chacha20poly1305.NewX(key)
	if err != nil {
		return nil, err
	}
	return &chachaClient{aead}, nil
}

// Encrypt encrypts plaintext with associated userID, returning hex-encoded ciphertext.
func (c *chachaClient) Encrypt(plaintext Plaintext, ownerID Owner) (Ciphertext, error) {
	nonce := make([]byte, c.aead.NonceSize())
	if _, err := cryptorand.Read(nonce); err != nil {
		return "", err
	}
	// Encrypt the message and append the ciphertext to the nonce, by using the
	// slice holding the nonce as the `dst` argument to `Seal()`.
	//
	// The result is a slice of bytes,
	// [ nonce : <encrypted message> : associated data ]
	encrypted := c.aead.Seal(nonce, nonce, []byte(plaintext), []byte(ownerID))
	return Ciphertext(hex.EncodeToString(encrypted)), nil
}

// Decrypt decrypts hex-encoded ciphertext with associated userID, returning plaintext.
func (c *chachaClient) Decrypt(ciphertext Ciphertext, ownerID Owner) (Plaintext, error) {
	ciphertextBytes, err := hex.DecodeString(string(ciphertext))
	if err != nil {
		return "", err
	}
	nonceSize := c.aead.NonceSize()
	nonce, encrypted := ciphertextBytes[:nonceSize], ciphertextBytes[nonceSize:]

	plaintext, err := c.aead.Open(nil, nonce, encrypted, []byte(ownerID))
	if err != nil {
		return "", err
	}
	return Plaintext(plaintext), nil
}

// GenerateRandomKey generates a random key that can be used by a [Client] to
// encrypt/decrypt messages. It uses a cryptographically-secure source of random
// bytes.
func GenerateRandomKey() (HexKey, error) {
	key := make([]byte, chacha20poly1305.KeySize)
	_, err := io.ReadFull(cryptorand.Reader, key)
	if err != nil {
		return "", err
	}
	return HexKey(hex.EncodeToString(key)), nil
}
