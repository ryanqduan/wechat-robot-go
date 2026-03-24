package crypto

import (
	"crypto/aes"
	"crypto/rand"
	"fmt"
)

// GenerateAESKey generates a random 16-byte AES-128 key.
func GenerateAESKey() ([]byte, error) {
	key := make([]byte, 16)
	_, err := rand.Read(key)
	return key, err
}

// PKCS7Pad pads data to a multiple of blockSize using PKCS7.
func PKCS7Pad(data []byte, blockSize int) []byte {
	padding := blockSize - len(data)%blockSize
	padtext := make([]byte, padding)
	for i := range padtext {
		padtext[i] = byte(padding)
	}
	return append(data, padtext...)
}

// PKCS7Unpad removes PKCS7 padding.
func PKCS7Unpad(data []byte) ([]byte, error) {
	if len(data) == 0 {
		return nil, fmt.Errorf("empty data")
	}
	padding := int(data[len(data)-1])
	if padding > len(data) || padding > aes.BlockSize || padding == 0 {
		return nil, fmt.Errorf("invalid padding size: %d", padding)
	}
	for i := len(data) - padding; i < len(data); i++ {
		if data[i] != byte(padding) {
			return nil, fmt.Errorf("invalid padding byte")
		}
	}
	return data[:len(data)-padding], nil
}

// EncryptAESECB encrypts data using AES-128-ECB mode with PKCS7 padding.
func EncryptAESECB(plaintext, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("create cipher: %w", err)
	}

	padded := PKCS7Pad(plaintext, aes.BlockSize)
	ciphertext := make([]byte, len(padded))

	for i := 0; i < len(padded); i += aes.BlockSize {
		block.Encrypt(ciphertext[i:i+aes.BlockSize], padded[i:i+aes.BlockSize])
	}

	return ciphertext, nil
}

// DecryptAESECB decrypts data using AES-128-ECB mode and removes PKCS7 padding.
func DecryptAESECB(ciphertext, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("create cipher: %w", err)
	}

	if len(ciphertext)%aes.BlockSize != 0 {
		return nil, fmt.Errorf("ciphertext length %d is not a multiple of block size %d", len(ciphertext), aes.BlockSize)
	}

	plaintext := make([]byte, len(ciphertext))
	for i := 0; i < len(ciphertext); i += aes.BlockSize {
		block.Decrypt(plaintext[i:i+aes.BlockSize], ciphertext[i:i+aes.BlockSize])
	}

	return PKCS7Unpad(plaintext)
}
