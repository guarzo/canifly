package persist

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/gob"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"os"
)

var key []byte

// Initialize sets up the encryption key. Key length must be 16, 24, or 32 bytes for AES-128, AES-192, or AES-256.
func Initialize(encryptionKey []byte) error {
	keyLength := len(encryptionKey)
	if keyLength != 16 && keyLength != 24 && keyLength != 32 {
		return fmt.Errorf("invalid key length: %d. Key must be 16, 24, or 32 bytes", keyLength)
	}
	key = encryptionKey
	return nil
}

// EncryptData encrypts the given data and writes it to the specified file.
func EncryptData(data interface{}, outputFile string) error {
	if !isKeyInitialized() {
		return errors.New("encryption key is not initialized")
	}

	outFile, err := os.OpenFile(outputFile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return err
	}
	defer outFile.Close()

	block, err := aes.NewCipher(key)
	if err != nil {
		return err
	}

	iv := make([]byte, aes.BlockSize)
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return err
	}

	if _, err := outFile.Write(iv); err != nil {
		return err
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	writer := &cipher.StreamWriter{S: stream, W: outFile}

	encoder := gob.NewEncoder(writer)
	if err := encoder.Encode(data); err != nil {
		return err
	}

	return nil
}

// DecryptData reads the encrypted data from the specified file, decrypts it, and populates the given data struct.
func DecryptData(inputFile string, data interface{}) error {
	if !isKeyInitialized() {
		return errors.New("decryption key is not initialized")
	}

	inFile, err := os.Open(inputFile)
	if err != nil {
		return err
	}
	defer inFile.Close()

	block, err := aes.NewCipher(key)
	if err != nil {
		return err
	}

	iv := make([]byte, aes.BlockSize)
	if _, err := io.ReadFull(inFile, iv); err != nil {
		return err
	}

	stream := cipher.NewCFBDecrypter(block, iv)
	reader := &cipher.StreamReader{S: stream, R: inFile}

	decoder := gob.NewDecoder(reader)
	if err := decoder.Decode(data); err != nil {
		return err
	}

	return nil
}

// isKeyInitialized checks if the encryption key is set.
func isKeyInitialized() bool {
	return len(key) > 0
}

func GenerateSecret() ([]byte, error) {
	newKey := make([]byte, 32)
	if _, err := rand.Read(newKey); err != nil {
		return nil, err
	}
	return newKey, nil
}

// EncryptString encrypts a plaintext string and returns a base64-encoded ciphertext string.
func EncryptString(plaintext string) (string, error) {
	if !isKeyInitialized() {
		return "", errors.New("encryption key is not initialized")
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	iv := make([]byte, aes.BlockSize)
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	ciphertext := make([]byte, len(plaintext))
	stream.XORKeyStream(ciphertext, []byte(plaintext))

	// Prepend IV to ciphertext
	finalData := append(iv, ciphertext...)
	// Return as base64
	return base64.StdEncoding.EncodeToString(finalData), nil
}

// DecryptString takes a base64-encoded ciphertext and returns the decrypted plaintext.
func DecryptString(ciphertextB64 string) (string, error) {
	if !isKeyInitialized() {
		return "", errors.New("decryption key is not initialized")
	}

	ciphertext, err := base64.StdEncoding.DecodeString(ciphertextB64)
	if err != nil {
		return "", err
	}

	if len(ciphertext) < aes.BlockSize {
		return "", errors.New("ciphertext too short")
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)
	plaintext := make([]byte, len(ciphertext))
	stream.XORKeyStream(plaintext, ciphertext)

	return string(plaintext), nil
}

func GenerateRandomString(lengthBytes int) (string, error) {
	b := make([]byte, lengthBytes)
	_, err := rand.Read(b)
	if err != nil {
		return "", fmt.Errorf("failed to generate random bytes: %w", err)
	}
	// Encode as hex string
	return hex.EncodeToString(b), nil
}

func ResetKeyForTest() {
	key = nil
}
