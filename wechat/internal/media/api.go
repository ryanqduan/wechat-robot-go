package media

import "context"

// APIClient is the interface for making API requests.
type APIClient interface {
	Post(ctx context.Context, path string, body, result interface{}) error
}
