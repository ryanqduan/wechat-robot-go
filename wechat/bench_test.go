package wechat

import (
	"testing"

	"github.com/ryanqduan/wechat-robot-go/wechat/internal/crypto"
	"github.com/ryanqduan/wechat-robot-go/wechat/internal/media"
)

func BenchmarkEncryptAESECB(b *testing.B) {
	key := []byte("0123456789abcdef") // 16 bytes
	data := make([]byte, 1024*1024)   // 1MB

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := crypto.EncryptAESECB(data, key)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkDecryptAESECB(b *testing.B) {
	key := []byte("0123456789abcdef")
	data := make([]byte, 1024*1024)
	encrypted, _ := crypto.EncryptAESECB(data, key)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := crypto.DecryptAESECB(encrypted, key)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkPKCS7Pad(b *testing.B) {
	data := make([]byte, 1000)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = crypto.PKCS7Pad(data, 16)
	}
}

func BenchmarkGenerateAESKey(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := crypto.GenerateAESKey()
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkBuildImageItem(b *testing.B) {
	client := NewClient("http://unused", nil, nil, "1.0.3")
	manager := media.NewMediaManager(client, client.HTTPClient(), nil)

	result := &media.UploadResult{
		AESKey:         "0123456789abcdef0123456789abcdef",
		FileKey:        "fedcba9876543210fedcba9876543210",
		EncryptedParam: "test-encrypted-param",
		FileSize:       12345,
		CipherSize:     12368,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = manager.BuildImageItem(result, 800, 600)
	}
}

func BenchmarkDecryptionOverhead(b *testing.B) {
	// This benchmark measures the decryption overhead
	key := []byte("0123456789abcdef")
	data := make([]byte, 1024)
	encrypted, _ := crypto.EncryptAESECB(data, key)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = crypto.DecryptAESECB(encrypted, key)
	}
}
