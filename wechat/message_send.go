package wechat

import (
	"context"
	"crypto/rand"
	"fmt"
	"math/big"
	"time"
)

// generateClientID generates a unique client ID for message tracking.
// Format: nanobot-{timestamp_ms}-{4_digit_random}
func generateClientID() string {
	randomPart := 1000
	n, err := rand.Int(rand.Reader, big.NewInt(9000))
	if err == nil {
		randomPart += int(n.Int64())
	}
	return fmt.Sprintf("nanobot-%d-%04d", time.Now().UnixMilli(), randomPart)
}

// SendText sends a text message to a user.
// contextToken must be provided from a received message for proper conversation linking.
func SendText(ctx context.Context, client *Client, toUserID, text, contextToken string) error {
	clientID := generateClientID()

	msg := &Message{
		FromUserID:   "", // Must be empty string (not omitted) per API spec
		ToUserID:     toUserID,
		ClientID:     clientID,
		MessageType:  MessageTypeBot,
		MessageState: MessageStateFinish,
		ContextToken: contextToken,
		ItemList: []MessageItem{
			{
				Type: ItemTypeText,
				TextItem: &TextItem{
					Text: text,
				},
			},
		},
	}

	req := &SendMessageRequest{
		Msg: msg,
		BaseInfo: &BaseInfo{
			ChannelVersion: client.channelVersion,
		},
	}

	var resp SendMessageResponse
	if err := client.Post(ctx, "/ilink/bot/sendmessage", req, &resp); err != nil {
		return err
	}

	if resp.Ret != 0 {
		return &APIError{Code: resp.ErrCode, Message: resp.ErrMsg}
	}

	return nil
}

// Reply sends a text reply to an incoming message, automatically using its context_token.
func Reply(ctx context.Context, client *Client, msg *Message, text string) error {
	return SendText(ctx, client, msg.FromUserID, text, msg.ContextToken)
}

// SendImage sends an image message to a user.
// The imageItem should contain the CDN-uploaded image information.
func SendImage(ctx context.Context, client *Client, toUserID, contextToken string, imageItem *ImageItem) error {
	msg := &Message{
		FromUserID:   "", // Must be empty string (not omitted) per API spec
		ToUserID:     toUserID,
		ClientID:     generateClientID(),
		MessageType:  MessageTypeBot,
		MessageState: MessageStateFinish,
		ContextToken: contextToken,
		ItemList: []MessageItem{
			{
				Type:      ItemTypeImage,
				ImageItem: imageItem,
			},
		},
	}

	req := &SendMessageRequest{
		Msg: msg,
		BaseInfo: &BaseInfo{
			ChannelVersion: client.channelVersion,
		},
	}

	var resp SendMessageResponse
	if err := client.Post(ctx, "/ilink/bot/sendmessage", req, &resp); err != nil {
		return err
	}

	if resp.Ret != 0 {
		return &APIError{Code: resp.ErrCode, Message: resp.ErrMsg}
	}

	return nil
}

// SendFile sends a file message to a user.
// The fileItem should contain the CDN-uploaded file information.
func SendFile(ctx context.Context, client *Client, toUserID, contextToken string, fileItem *FileItem) error {
	msg := &Message{
		FromUserID:   "", // Must be empty string (not omitted) per API spec
		ToUserID:     toUserID,
		ClientID:     generateClientID(),
		MessageType:  MessageTypeBot,
		MessageState: MessageStateFinish,
		ContextToken: contextToken,
		ItemList: []MessageItem{
			{
				Type:     ItemTypeFile,
				FileItem: fileItem,
			},
		},
	}

	req := &SendMessageRequest{
		Msg: msg,
		BaseInfo: &BaseInfo{
			ChannelVersion: client.channelVersion,
		},
	}

	var resp SendMessageResponse
	if err := client.Post(ctx, "/ilink/bot/sendmessage", req, &resp); err != nil {
		return err
	}

	if resp.Ret != 0 {
		return &APIError{Code: resp.ErrCode, Message: resp.ErrMsg}
	}

	return nil
}

// SendMessage sends a custom message with multiple items.
func SendMessage(ctx context.Context, client *Client, toUserID, contextToken string, items []MessageItem) error {
	msg := &Message{
		FromUserID:   "", // Must be empty string (not omitted) per API spec
		ToUserID:     toUserID,
		ClientID:     generateClientID(),
		MessageType:  MessageTypeBot,
		MessageState: MessageStateFinish,
		ContextToken: contextToken,
		ItemList:     items,
	}

	req := &SendMessageRequest{
		Msg: msg,
		BaseInfo: &BaseInfo{
			ChannelVersion: client.channelVersion,
		},
	}

	var resp SendMessageResponse
	if err := client.Post(ctx, "/ilink/bot/sendmessage", req, &resp); err != nil {
		return err
	}

	if resp.Ret != 0 {
		return &APIError{Code: resp.ErrCode, Message: resp.ErrMsg}
	}

	return nil
}

// SendVoice sends a voice message to a user.
// The voiceItem should contain the CDN-uploaded voice information.
func SendVoice(ctx context.Context, client *Client, toUserID, contextToken string, voiceItem *VoiceItem) error {
	msg := &Message{
		FromUserID:   "", // Must be empty string (not omitted) per API spec
		ToUserID:     toUserID,
		ClientID:     generateClientID(),
		MessageType:  MessageTypeBot,
		MessageState: MessageStateFinish,
		ContextToken: contextToken,
		ItemList: []MessageItem{
			{
				Type:      ItemTypeVoice,
				VoiceItem: voiceItem,
			},
		},
	}

	req := &SendMessageRequest{
		Msg: msg,
		BaseInfo: &BaseInfo{
			ChannelVersion: client.channelVersion,
		},
	}

	var resp SendMessageResponse
	if err := client.Post(ctx, "/ilink/bot/sendmessage", req, &resp); err != nil {
		return err
	}

	if resp.Ret != 0 {
		return &APIError{Code: resp.ErrCode, Message: resp.ErrMsg}
	}

	return nil
}

// SendVideo sends a video message to a user.
// The videoItem should contain the CDN-uploaded video information.
func SendVideo(ctx context.Context, client *Client, toUserID, contextToken string, videoItem *VideoItem) error {
	msg := &Message{
		FromUserID:   "", // Must be empty string (not omitted) per API spec
		ToUserID:     toUserID,
		ClientID:     generateClientID(),
		MessageType:  MessageTypeBot,
		MessageState: MessageStateFinish,
		ContextToken: contextToken,
		ItemList: []MessageItem{
			{
				Type:      ItemTypeVideo,
				VideoItem: videoItem,
			},
		},
	}

	req := &SendMessageRequest{
		Msg: msg,
		BaseInfo: &BaseInfo{
			ChannelVersion: client.channelVersion,
		},
	}

	var resp SendMessageResponse
	if err := client.Post(ctx, "/ilink/bot/sendmessage", req, &resp); err != nil {
		return err
	}

	if resp.Ret != 0 {
		return &APIError{Code: resp.ErrCode, Message: resp.ErrMsg}
	}

	return nil
}

// ReplyWithMedia sends a rich media reply with both text and media items.
func ReplyWithMedia(ctx context.Context, client *Client, msg *Message, text string, mediaItems []MessageItem) error {
	items := []MessageItem{
		{
			Type: ItemTypeText,
			TextItem: &TextItem{
				Text: text,
			},
		},
	}
	items = append(items, mediaItems...)

	return SendMessage(ctx, client, msg.FromUserID, msg.ContextToken, items)
}
