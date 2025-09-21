package tools

import "context"

type Storage interface {
	Store(ctx context.Context, image string) (string, error)
}