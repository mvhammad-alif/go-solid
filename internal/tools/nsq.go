package tools

import "context"

type Publisher interface {
	Publish(ctx context.Context, topic, message string) (error)
}