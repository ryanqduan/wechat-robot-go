package wechat

import (
	"errors"
	"fmt"
)

// APIError represents an error returned by the iLink Bot API.
type APIError struct {
	Ret     int    `json:"ret"`
	ErrCode int    `json:"errcode"`
	Message string `json:"errmsg"`
	// Code is retained for callers that inspect a single error code.
	// SDK responses use ErrCode when nonzero, otherwise Ret.
	Code int `json:"-"`
}

func (e *APIError) Error() string {
	// Preserve diagnostics for callers constructing the legacy Code-only form.
	if e.Ret == 0 && e.ErrCode == 0 && e.Code != 0 {
		return fmt.Sprintf("ilink api error: code=%d, msg=%s", e.Code, e.Message)
	}
	return fmt.Sprintf("ilink api error: ret=%d, errcode=%d, msg=%s", e.Ret, e.ErrCode, e.Message)
}

// apiResponseError preserves both response codes. Missing JSON codes decode to zero;
// successful responses may omit them (for example, sendmessage's message_id response).
func apiResponseError(ret, errCode int, message string) error {
	if ret == 0 && errCode == 0 {
		return nil
	}
	code := errCode
	if code == 0 {
		code = ret
	}
	return &APIError{Ret: ret, ErrCode: errCode, Code: code, Message: message}
}

// IsSessionExpired reports whether either response code indicates expiration (-14).
func IsSessionExpired(err error) bool {
	var apiErr *APIError
	if errors.As(err, &apiErr) {
		return apiErr.Ret == -14 || apiErr.ErrCode == -14 || apiErr.Code == -14
	}
	return false
}

// Sentinel errors
var (
	ErrSessionExpired = &APIError{ErrCode: -14, Code: -14, Message: "session expired"}
	ErrNotLoggedIn    = errors.New("wechat: not logged in")
	ErrQRCodeExpired  = errors.New("wechat: qr code expired")
	ErrPollerStopped  = errors.New("wechat: poller stopped")
)
