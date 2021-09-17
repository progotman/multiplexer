package requester

import "context"

type Interface interface {
	GetContent(ctx context.Context, url string) (string, error)
}
