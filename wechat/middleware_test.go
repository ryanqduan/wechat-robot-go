package wechat

import (
	"context"
	"errors"
	"log/slog"
	"os"
	"sync/atomic"
	"testing"
)

func TestChain(t *testing.T) {
	var order []string

	m1 := func(next MessageHandler) MessageHandler {
		return func(ctx context.Context, msg *Message) error {
			order = append(order, "m1-before")
			err := next(ctx, msg)
			order = append(order, "m1-after")
			return err
		}
	}
	m2 := func(next MessageHandler) MessageHandler {
		return func(ctx context.Context, msg *Message) error {
			order = append(order, "m2-before")
			err := next(ctx, msg)
			order = append(order, "m2-after")
			return err
		}
	}

	handler := func(ctx context.Context, msg *Message) error {
		order = append(order, "handler")
		return nil
	}

	chained := Chain(m1, m2)(handler)
	_ = chained(context.Background(), &Message{})

	expected := []string{"m1-before", "m2-before", "handler", "m2-after", "m1-after"}
	if len(order) != len(expected) {
		t.Fatalf("got %v, want %v", order, expected)
	}
	for i, v := range expected {
		if order[i] != v {
			t.Errorf("order[%d] = %q, want %q", i, order[i], v)
		}
	}
}

func TestWithRecovery(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelError}))
	middleware := WithRecovery(logger)

	handler := middleware(func(ctx context.Context, msg *Message) error {
		panic("test panic in middleware")
	})

	err := handler(context.Background(), &Message{FromUserID: "user1"})
	if err == nil {
		t.Error("expected error from recovered panic")
	}
	if err != nil && err.Error() != "handler panic: test panic in middleware" {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestWithLogging(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelDebug}))
	var called int32

	middleware := WithLogging(logger)
	handler := middleware(func(ctx context.Context, msg *Message) error {
		atomic.AddInt32(&called, 1)
		return nil
	})

	_ = handler(context.Background(), &Message{FromUserID: "user1", ItemList: []MessageItem{{Type: ItemTypeText}}})

	if atomic.LoadInt32(&called) != 1 {
		t.Error("handler was not called")
	}
}

func TestWithLogging_Error(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelDebug}))

	middleware := WithLogging(logger)
	handler := middleware(func(ctx context.Context, msg *Message) error {
		return errors.New("test error")
	})

	err := handler(context.Background(), &Message{FromUserID: "user1"})
	if err == nil || err.Error() != "test error" {
		t.Errorf("expected 'test error', got %v", err)
	}
}
