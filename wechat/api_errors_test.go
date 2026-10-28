package wechat

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"
)

// Exercise the HTTP decoding and public API paths, including the response that
// previously produced "code=0, msg=prepare failed" and message_id-only success.
func TestAPIResponseErrors(t *testing.T) {
	ctx := context.Background()
	typingManager := func(c *Client) *TypingManager {
		store := NewMemoryContextTokenStore()
		if err := store.Save("user", "context"); err != nil {
			t.Fatal(err)
		}
		return NewTypingManager(c, store, slog.Default())
	}
	operations := []struct {
		name string
		path string
		call func(*Client) error
	}{
		{"SendText", "sendmessage", func(c *Client) error { return SendText(ctx, c, "user", "text", "context") }},
		{"SendImage", "sendmessage", func(c *Client) error { return SendImage(ctx, c, "user", "context", &ImageItem{}) }},
		{"SendFile", "sendmessage", func(c *Client) error { return SendFile(ctx, c, "user", "context", &FileItem{}) }},
		{"SendVoice", "sendmessage", func(c *Client) error { return SendVoice(ctx, c, "user", "context", &VoiceItem{}) }},
		{"SendVideo", "sendmessage", func(c *Client) error { return SendVideo(ctx, c, "user", "context", &VideoItem{}) }},
		{"SendMessage", "sendmessage", func(c *Client) error {
			return SendMessage(ctx, c, "user", "context", []MessageItem{{Type: ItemTypeText, TextItem: &TextItem{Text: "text"}}})
		}},
		{"SendImageWithItem", "sendmessage", func(c *Client) error { return SendImageWithItem(ctx, c, "user", "context", &ImageItem{}) }},
		{"SendFileWithItem", "sendmessage", func(c *Client) error { return SendFileWithItem(ctx, c, "user", "context", &FileItem{}) }},
		{"SendVoiceWithItem", "sendmessage", func(c *Client) error { return SendVoiceWithItem(ctx, c, "user", "context", &VoiceItem{}) }},
		{"SendVideoWithItem", "sendmessage", func(c *Client) error { return SendVideoWithItem(ctx, c, "user", "context", &VideoItem{}) }},
		{"Poll", "getupdates", func(c *Client) error {
			_, err := NewPoller(c, nil, slog.Default(), "1.0.3").poll(ctx)
			return err
		}},
		{"ValidateCredentials", "getupdates", func(c *Client) error {
			valid, err := NewAuth(c, nil, slog.Default()).ValidateCredentials(ctx)
			if valid != (err == nil) {
				t.Errorf("valid=%v inconsistent with error=%v", valid, err)
			}
			return err
		}},
		{"GetConfig", "getconfig", func(c *Client) error {
			_, err := typingManager(c).GetConfig(ctx, "user")
			return err
		}},
		{"SendTyping", "sendtyping", func(c *Client) error { return typingManager(c).SendTyping(ctx, "user") }},
		{"StopTyping", "sendtyping", func(c *Client) error { return typingManager(c).StopTyping(ctx, "user") }},
	}
	responses := []struct {
		name    string
		body    string
		ret     int
		errCode int
		code    int
		message string
	}{
		{"ret_only", `{"ret":-2,"errmsg":"prepare failed"}`, -2, 0, -2, "prepare failed"},
		{"explicit_zero_errcode", `{"ret":-2,"errcode":0,"errmsg":"prepare failed"}`, -2, 0, -2, "prepare failed"},
		{"errcode_only", `{"errcode":-14,"errmsg":"expired"}`, 0, -14, -14, "expired"},
		{"zero_ret_nonzero_errcode", `{"ret":0,"errcode":-14,"errmsg":"expired"}`, 0, -14, -14, "expired"},
		{"both_codes", `{"ret":-1,"errcode":-1001,"errmsg":"send failed"}`, -1, -1001, -1001, "send failed"},
		{"expired_ret_other_errcode", `{"ret":-14,"errcode":-1,"errmsg":"expired"}`, -14, -1, -1, "expired"},
		{"no_message", `{"ret":-2}`, -2, 0, -2, ""},
		{"success", `{"ret":0,"errcode":0,"errmsg":"ok"}`, 0, 0, 0, "ok"},
		{"message_id_only", `{"message_id":"12345678901234567890"}`, 0, 0, 0, ""},
	}
	for _, op := range operations {
		t.Run(op.name, func(t *testing.T) {
			for _, response := range responses {
				t.Run(response.name, func(t *testing.T) {
					server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
						w.Header().Set("Content-Type", "application/json")
						if op.path == "sendtyping" && r.URL.Path == "/ilink/bot/getconfig" {
							fmt.Fprint(w, `{"ret":0,"typing_ticket":"ticket"}`)
							return
						}
						if r.URL.Path != "/ilink/bot/"+op.path {
							t.Errorf("unexpected path: %s", r.URL.Path)
						}
						fmt.Fprint(w, response.body)
					}))
					defer server.Close()
					client := NewClient(server.URL, server.Client(), slog.Default(), "1.0.3")
					err := op.call(client)
					if response.code == 0 {
						if err != nil {
							t.Fatalf("successful response rejected: %v", err)
						}
						return
					}
					var apiErr *APIError
					if !errors.As(err, &apiErr) {
						t.Fatalf("expected APIError, got %v", err)
					}
					if apiErr.Ret != response.ret || apiErr.ErrCode != response.errCode || apiErr.Code != response.code || apiErr.Message != response.message {
						t.Fatalf("incorrect error fields: %#v", apiErr)
					}
					want := fmt.Sprintf("ilink api error: ret=%d, errcode=%d, msg=%s", response.ret, response.errCode, response.message)
					if apiErr.Error() != want {
						t.Errorf("error=%q, want %q", apiErr.Error(), want)
					}
					if got, want := IsSessionExpired(err), response.ret == -14 || response.errCode == -14; got != want {
						t.Errorf("IsSessionExpired=%v, want %v", got, want)
					}
				})
			}
		})
	}
}

func TestAPIErrorJSONPreservesResponseCodes(t *testing.T) {
	err := apiResponseError(-2, 0, "prepare failed")
	data, marshalErr := json.Marshal(err)
	if marshalErr != nil {
		t.Fatal(marshalErr)
	}
	var got struct {
		Ret     int    `json:"ret"`
		ErrCode int    `json:"errcode"`
		Message string `json:"errmsg"`
	}
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatal(err)
	}
	if got.Ret != -2 || got.ErrCode != 0 || got.Message != "prepare failed" {
		t.Fatalf("response codes changed during JSON encoding: %s", data)
	}
}

func TestIsSessionExpiredCompatibility(t *testing.T) {
	for _, err := range []error{
		ErrSessionExpired,
		&APIError{Code: -14, Message: "legacy"},
		&APIError{Ret: -14, ErrCode: -1},
		&APIError{Ret: -1, ErrCode: -14},
	} {
		if !IsSessionExpired(fmt.Errorf("wrapped: %w", err)) {
			t.Errorf("expiration not detected: %v", err)
		}
	}
	for _, err := range []error{nil, errors.New("network error"), &APIError{Ret: -2}} {
		if IsSessionExpired(err) {
			t.Errorf("unexpected expiration: %v", err)
		}
	}
}
