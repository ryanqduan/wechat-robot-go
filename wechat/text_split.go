package wechat

import (
	"context"
	"fmt"

	"github.com/SpellingDragon/wechat-robot-go/wechat/internal/text"
)

// SendLongText sends a long text message by automatically splitting it into multiple messages.
// It returns the number of messages sent and any error encountered.
func SendLongText(ctx context.Context, client *Client, toUserID, textContent, contextToken string) (int, error) {
	chunks := text.SplitText(textContent, text.DefaultMaxTextLength)

	for i, chunk := range chunks {
		if err := SendText(ctx, client, toUserID, chunk, contextToken); err != nil {
			return i, fmt.Errorf("send chunk %d: %w", i+1, err)
		}
	}

	return len(chunks), nil
}
