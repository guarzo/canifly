package persist_test

import (
	"path/filepath"
	"testing"

	"github.com/guarzo/canifly/internal/persist"
	"github.com/stretchr/testify/assert"
)

func TestInitialize(t *testing.T) {
	// Invalid key length
	err := persist.Initialize([]byte("short"))
	assert.Error(t, err, "Should fail with invalid key length")

	// Valid 16-byte key
	key16 := make([]byte, 16)
	err = persist.Initialize(key16)
	assert.NoError(t, err, "Should succeed with 16-byte key")

	// Valid 24-byte key
	key24 := make([]byte, 24)
	err = persist.Initialize(key24)
	assert.NoError(t, err, "Should succeed with 24-byte key")

	// Valid 32-byte key
	key32 := make([]byte, 32)
	err = persist.Initialize(key32)
	assert.NoError(t, err, "Should succeed with 32-byte key")
}

func TestEncryptDecryptData(t *testing.T) {
	// Initialize with a 32-byte key
	key := make([]byte, 32)
	err := persist.Initialize(key)
	assert.NoError(t, err)

	type TestData struct {
		Name  string
		Value int
	}

	inputData := TestData{
		Name:  "example",
		Value: 42,
	}

	tempDir := t.TempDir()
	filePath := filepath.Join(tempDir, "encrypted_data.bin")

	// Encrypt data
	err = persist.EncryptData(inputData, filePath)
	assert.NoError(t, err, "Encryption should succeed")

	// Decrypt data
	var outputData TestData
	err = persist.DecryptData(filePath, &outputData)
	assert.NoError(t, err, "Decryption should succeed")

	assert.Equal(t, inputData.Name, outputData.Name)
	assert.Equal(t, inputData.Value, outputData.Value)
}

func TestEncryptDecryptData_NoKey(t *testing.T) {
	persist.ResetKeyForTest()
	type TestData struct{ Foo string }
	data := TestData{Foo: "bar"}

	tempDir := t.TempDir()
	filePath := filepath.Join(tempDir, "no_key.bin")

	// Attempting to EncryptData without initializing key should fail
	err := persist.EncryptData(data, filePath)
	assert.Error(t, err)

	err = persist.DecryptData(filePath, &data)
	assert.Error(t, err)
}

func TestEncryptDecryptString(t *testing.T) {
	key := make([]byte, 32)
	err := persist.Initialize(key)
	assert.NoError(t, err)

	plaintext := "Hello, World!"
	encrypted, err := persist.EncryptString(plaintext)
	assert.NoError(t, err)
	assert.NotEmpty(t, encrypted)

	decrypted, err := persist.DecryptString(encrypted)
	assert.NoError(t, err)
	assert.Equal(t, plaintext, decrypted)
}

func TestEncryptDecryptString_NoKey(t *testing.T) {
	persist.ResetKeyForTest()
	_, err := persist.EncryptString("test")
	assert.Error(t, err)

	_, err = persist.DecryptString("somecipher")
	assert.Error(t, err)
}

func TestGenerateSecret(t *testing.T) {
	secret, err := persist.GenerateSecret()
	assert.NoError(t, err)
	assert.Len(t, secret, 32, "GenerateSecret should return 32 bytes")

	// Try initializing with the generated secret
	err = persist.Initialize(secret)
	assert.NoError(t, err)
}

func TestGenerateRandomString(t *testing.T) {
	str, err := persist.GenerateRandomString(16)
	assert.NoError(t, err)
	assert.Len(t, str, 32, "Hex encoded length should be double the byte length")

	// Try another length
	str2, err := persist.GenerateRandomString(10)
	assert.NoError(t, err)
	assert.Len(t, str2, 20)
}

func TestDecryptData_FileNotExist(t *testing.T) {
	key := make([]byte, 32)
	err := persist.Initialize(key)
	assert.NoError(t, err)

	// File doesn't exist
	var data interface{}
	err = persist.DecryptData("non_existent.bin", &data)
	assert.Error(t, err)
}

func TestEncryptData_WriteError(t *testing.T) {
	key := make([]byte, 32)
	err := persist.Initialize(key)
	assert.NoError(t, err)

	// Invalid path (unwritable directory)
	err = persist.EncryptData("data", "/root/forbidden.bin")
	assert.Error(t, err)
}
