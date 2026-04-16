package media

import (
	"encoding/base64"
	"fmt"

	"github.com/ryanqduan/wechat-robot-go/wechat/internal/model"
)

// BuildImageItem creates a MessageItem for sending an image.
// This follows the official Weixin API format with CDNMedia structure.
// Note: AES key format for outbound: base64(hex_string), NOT base64(raw_bytes).
// Official plugin: Buffer.from(aeskey_hex).toString("base64") = base64 encode hex string as UTF-8.
func (m *MediaManager) BuildImageItem(result *UploadResult, width, height int) model.MessageItem {
	// Per official plugin: aes_key = base64(hex_string)
	// NOT base64(raw_bytes)! This is counter-intuitive but critical.
	aesKeyBase64 := base64.StdEncoding.EncodeToString([]byte(result.AESKey))

	return model.MessageItem{
		Type: model.ItemTypeImage,
		ImageItem: &model.ImageItem{
			Media: &model.CDNMedia{
				EncryptQueryParam: result.EncryptedParam,
				AESKey:            aesKeyBase64, // base64(hex_string) per official plugin
				EncryptType:       1,            // encrypted with bundled thumbnail info
			},
			MidSize: result.CipherSize, // ciphertext file size
		},
	}
}

// BuildImageItemPtr creates an ImageItem for sending an image.
func (m *MediaManager) BuildImageItemPtr(result *UploadResult, width, height int) *model.ImageItem {
	// Per official plugin: aes_key = base64(hex_string)
	aesKeyBase64 := base64.StdEncoding.EncodeToString([]byte(result.AESKey))

	return &model.ImageItem{
		Media: &model.CDNMedia{
			EncryptQueryParam: result.EncryptedParam,
			AESKey:            aesKeyBase64,
			EncryptType:       1,
		},
		MidSize: result.CipherSize,
	}
}

// BuildFileItem creates a MessageItem for sending a file.
func (m *MediaManager) BuildFileItem(result *UploadResult, fileName string) model.MessageItem {
	// Per official plugin: aes_key = base64(hex_string)
	aesKeyBase64 := base64.StdEncoding.EncodeToString([]byte(result.AESKey))

	return model.MessageItem{
		Type: model.ItemTypeFile,
		FileItem: &model.FileItem{
			Media: &model.CDNMedia{
				EncryptQueryParam: result.EncryptedParam,
				AESKey:            aesKeyBase64,
				EncryptType:       1,
			},
			FileName: fileName,
			Length:   fmt.Sprintf("%d", result.FileSize), // plaintext size as string
		},
	}
}

// BuildFileItemPtr creates a FileItem for sending a file.
func (m *MediaManager) BuildFileItemPtr(result *UploadResult, fileName string) *model.FileItem {
	// Per official plugin: aes_key = base64(hex_string)
	aesKeyBase64 := base64.StdEncoding.EncodeToString([]byte(result.AESKey))

	return &model.FileItem{
		Media: &model.CDNMedia{
			EncryptQueryParam: result.EncryptedParam,
			AESKey:            aesKeyBase64,
			EncryptType:       1,
		},
		FileName: fileName,
		Length:   fmt.Sprintf("%d", result.FileSize),
	}
}

// BuildVideoItem creates a MessageItem for sending a video.
func (m *MediaManager) BuildVideoItem(result *UploadResult, width, height, duration int) model.MessageItem {
	// Per official plugin: aes_key = base64(hex_string)
	aesKeyBase64 := base64.StdEncoding.EncodeToString([]byte(result.AESKey))

	return model.MessageItem{
		Type: model.ItemTypeVideo,
		VideoItem: &model.VideoItem{
			Media: &model.CDNMedia{
				EncryptQueryParam: result.EncryptedParam,
				AESKey:            aesKeyBase64,
				EncryptType:       1,
			},
			VideoSize:   result.FileSize,
			PlayLength:  duration,
			ThumbWidth:  width,
			ThumbHeight: height,
		},
	}
}

// BuildVideoItemPtr creates a VideoItem for sending a video.
func (m *MediaManager) BuildVideoItemPtr(result *UploadResult, width, height, duration int) *model.VideoItem {
	// Per official plugin: aes_key = base64(hex_string)
	aesKeyBase64 := base64.StdEncoding.EncodeToString([]byte(result.AESKey))

	return &model.VideoItem{
		Media: &model.CDNMedia{
			EncryptQueryParam: result.EncryptedParam,
			AESKey:            aesKeyBase64,
			EncryptType:       1,
		},
		VideoSize:   result.FileSize,
		PlayLength:  duration,
		ThumbWidth:  width,
		ThumbHeight: height,
	}
}

// BuildVoiceItem creates a MessageItem for sending a voice message.
func (m *MediaManager) BuildVoiceItem(result *UploadResult, duration int) model.MessageItem {
	// Per official plugin: aes_key = base64(hex_string)
	aesKeyBase64 := base64.StdEncoding.EncodeToString([]byte(result.AESKey))

	return model.MessageItem{
		Type: model.ItemTypeVoice,
		VoiceItem: &model.VoiceItem{
			Media: &model.CDNMedia{
				EncryptQueryParam: result.EncryptedParam,
				AESKey:            aesKeyBase64,
				EncryptType:       1,
			},
			Duration: duration, // in milliseconds
		},
	}
}

// BuildVoiceItemPtr creates a VoiceItem for sending a voice message.
func (m *MediaManager) BuildVoiceItemPtr(result *UploadResult, duration int) *model.VoiceItem {
	// Per official plugin: aes_key = base64(hex_string)
	aesKeyBase64 := base64.StdEncoding.EncodeToString([]byte(result.AESKey))

	return &model.VoiceItem{
		Media: &model.CDNMedia{
			EncryptQueryParam: result.EncryptedParam,
			AESKey:            aesKeyBase64,
			EncryptType:       1,
		},
		Duration: duration,
	}
}
