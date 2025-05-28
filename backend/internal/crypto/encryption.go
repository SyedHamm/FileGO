package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"io"
)

// Constants
const (
	KeySize = 32 // 256-bit key
)

// GenerateRandomKey generates a random key for encryption
func GenerateRandomKey() ([]byte, error) {
	key := make([]byte, KeySize)
	_, err := io.ReadFull(rand.Reader, key)
	if err != nil {
		return nil, err
	}
	return key, nil
}

// EncryptFile encrypts a file using AES-GCM
func EncryptFile(src io.Reader, dst io.Writer, key []byte) error {
	// Create a new cipher block from the key
	block, err := aes.NewCipher(key)
	if err != nil {
		return err
	}

	// Create a new GCM
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return err
	}

	// Create a random nonce
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return err
	}

	// Write the nonce to the output file
	if _, err := dst.Write(nonce); err != nil {
		return err
	}

	// Create a writer that will encrypt and write to the destination
	encryptWriter := &cipher.StreamWriter{
		S: cipher.NewOFB(block, nonce),
		W: dst,
	}

	// Copy the input file to the encrypted output writer
	if _, err := io.Copy(encryptWriter, src); err != nil {
		return err
	}

	return nil
}

// DecryptFile decrypts a file using AES-GCM
func DecryptFile(src io.Reader, dst io.Writer, key []byte) error {
	// Create a new cipher block from the key
	block, err := aes.NewCipher(key)
	if err != nil {
		return err
	}

	// Create a new GCM
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return err
	}

	// Read the nonce from the encrypted file
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(src, nonce); err != nil {
		return err
	}

	// Create a reader that will decrypt from the source
	decryptReader := &cipher.StreamReader{
		S: cipher.NewOFB(block, nonce),
		R: src,
	}

	// Copy the decrypted input to the output file
	if _, err := io.Copy(dst, decryptReader); err != nil {
		return err
	}

	return nil
}

// KeyToString converts a key to a hex string
func KeyToString(key []byte) string {
	return hex.EncodeToString(key)
}

// StringToKey converts a hex string to a key
func StringToKey(s string) ([]byte, error) {
	return hex.DecodeString(s)
}
