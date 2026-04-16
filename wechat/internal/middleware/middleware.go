package middleware

import (
	"context"
	"fmt"
	"log/slog"
	"runtime/debug"

	"github.com/ryanqduan/wechat-robot-go/wechat/internal/model"
)

// Middleware wraps a MessageHandler to add cross-cutting concerns.
// Middlewares are applied in the order they are registered:
// the first middleware is the outermost wrapper.
type Middleware func(model.MessageHandler) model.MessageHandler

// Chain composes multiple middlewares into a single Middleware.
// Middlewares are applied in the order provided:
//
//	Chain(A, B, C)(handler) == A(B(C(handler)))
func Chain(middlewares ...Middleware) Middleware {
	return func(next model.MessageHandler) model.MessageHandler {
		for i := len(middlewares) - 1; i >= 0; i-- {
			next = middlewares[i](next)
		}
		return next
	}
}

// WithRecovery returns a middleware that recovers from panics in the handler.
// Panics are logged with stack trace and do not propagate.
func WithRecovery(logger *slog.Logger) Middleware {
	return func(next model.MessageHandler) model.MessageHandler {
		return func(ctx context.Context, msg *model.Message) (err error) {
			defer func() {
				if r := recover(); r != nil {
					stack := string(debug.Stack())
					logger.Error("handler panic recovered",
						slog.String("from_user_id", msg.FromUserID),
						slog.Any("panic", r),
						slog.String("stack", stack),
					)
					err = fmt.Errorf("handler panic: %v", r)
				}
			}()
			return next(ctx, msg)
		}
	}
}

// WithLogging returns a middleware that logs incoming messages and handler results.
func WithLogging(logger *slog.Logger) Middleware {
	return func(next model.MessageHandler) model.MessageHandler {
		return func(ctx context.Context, msg *model.Message) error {
			logger.Info("handling message",
				slog.String("from_user_id", msg.FromUserID),
				slog.Int("item_count", len(msg.ItemList)),
			)

			err := next(ctx, msg)

			if err != nil {
				logger.Warn("handler returned error",
					slog.String("from_user_id", msg.FromUserID),
					slog.String("error", err.Error()),
				)
			}

			return err
		}
	}
}
