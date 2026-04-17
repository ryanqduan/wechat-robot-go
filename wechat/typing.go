package wechat

import (
	"context"
	"log/slog"
	"sync"
	"time"
)

const (
	// TypingStatusStart indicates typing has started.
	TypingStatusStart = 1
	// TypingStatusStop indicates typing has stopped.
	TypingStatusStop = 2

	// ticketCacheDuration is how long to cache the typing_ticket.
	ticketCacheDuration = 24 * time.Hour
)

// TypingManager manages "typing" status indicators.
type TypingManager struct {
	client        *Client
	contextTokens ContextTokenStore
	logger        *slog.Logger
	typingTickets map[string]typingTicketCacheEntry // user_id -> cached typing ticket
	mu            sync.RWMutex
}

type typingTicketCacheEntry struct {
	ticket string
	expiry time.Time
}

// SendTypingRequest is the request body for POST /ilink/bot/sendtyping.
type SendTypingRequest struct {
	ILinkUserID  string    `json:"ilink_user_id"`
	TypingTicket string    `json:"typing_ticket"`
	Status       int       `json:"status"` // 1=typing, 2=cancel
	BaseInfo     *BaseInfo `json:"base_info,omitempty"`
}

type GetConfigRequest struct {
	ILinkUserID  string    `json:"ilink_user_id"`
	ContextToken string    `json:"context_token"`
	BaseInfo     *BaseInfo `json:"base_info,omitempty"`
}

// SendTypingResponse is the response from POST /ilink/bot/sendtyping.
type SendTypingResponse struct {
	Ret     int    `json:"ret"`
	ErrCode int    `json:"errcode,omitempty"`
	ErrMsg  string `json:"errmsg,omitempty"`
}

// NewTypingManager creates a new TypingManager instance.
func NewTypingManager(client *Client, contextTokens ContextTokenStore, logger *slog.Logger) *TypingManager {
	return &TypingManager{
		client:        client,
		contextTokens: contextTokens,
		logger:        logger,
		typingTickets: make(map[string]typingTicketCacheEntry),
	}
}

// GetConfig fetches the typing_ticket from the server.
// Caches the ticket for 24 hours.
func (tm *TypingManager) GetConfig(ctx context.Context, toUserID string) (string, error) {
	// Check cache first
	tm.mu.RLock()
	if cached, ok := tm.typingTickets[toUserID]; ok && cached.ticket != "" && time.Now().Before(cached.expiry) {
		ticket := cached.ticket
		tm.mu.RUnlock()
		return ticket, nil
	}
	tm.mu.RUnlock()

	// Query from contextTokens
	contextToken, err := tm.contextTokens.Load(toUserID)
	if err != nil {
		return "", err
	}
	if contextToken == "" {
		return "", &APIError{Code: 40001, Message: "no context token found for user"}
	}
	req := &GetConfigRequest{
		ILinkUserID:  toUserID,
		ContextToken: contextToken,
		BaseInfo: &BaseInfo{
			ChannelVersion: tm.client.channelVersion,
		},
	}

	// Fetch from server
	var resp GetConfigResponse
	if err := tm.client.Post(ctx, "/ilink/bot/getconfig", req, &resp); err != nil {
		return "", err
	}

	if resp.Ret != 0 {
		return "", &APIError{Code: resp.Ret, Message: resp.ErrMsg}
	}

	// Cache the ticket
	tm.mu.Lock()
	tm.typingTickets[toUserID] = typingTicketCacheEntry{
		ticket: resp.TypingTicket,
		expiry: time.Now().Add(ticketCacheDuration),
	}
	tm.mu.Unlock()

	return resp.TypingTicket, nil
}

// SendTyping sends a "typing" indicator to the user.
func (tm *TypingManager) SendTyping(ctx context.Context, toUserID string) error {
	ticket, err := tm.GetConfig(ctx, toUserID)
	if err != nil {
		return err
	}

	req := &SendTypingRequest{
		ILinkUserID:  toUserID,
		TypingTicket: ticket,
		Status:       TypingStatusStart,
		BaseInfo: &BaseInfo{
			ChannelVersion: tm.client.channelVersion,
		},
	}

	var resp SendTypingResponse
	if err := tm.client.Post(ctx, "/ilink/bot/sendtyping", req, &resp); err != nil {
		return err
	}

	if resp.Ret != 0 {
		return &APIError{Code: resp.Ret, Message: resp.ErrMsg}
	}

	return nil
}

// StopTyping cancels the "typing" indicator.
func (tm *TypingManager) StopTyping(ctx context.Context, toUserID string) error {
	ticket, err := tm.GetConfig(ctx, toUserID)
	if err != nil {
		return err
	}

	req := &SendTypingRequest{
		ILinkUserID:  toUserID,
		TypingTicket: ticket,
		Status:       TypingStatusStop,
		BaseInfo: &BaseInfo{
			ChannelVersion: tm.client.channelVersion,
		},
	}

	var resp SendTypingResponse
	if err := tm.client.Post(ctx, "/ilink/bot/sendtyping", req, &resp); err != nil {
		return err
	}

	if resp.Ret != 0 {
		return &APIError{Code: resp.Ret, Message: resp.ErrMsg}
	}

	return nil
}

// ClearCache clears the cached typing ticket, forcing a refresh on next call.
func (tm *TypingManager) ClearCache() {
	tm.mu.Lock()
	tm.typingTickets = make(map[string]typingTicketCacheEntry)
	tm.mu.Unlock()
}
