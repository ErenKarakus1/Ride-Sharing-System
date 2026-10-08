package events

import "context"

type Event struct {
	Type string
	Key  string
	Data any
}

type Publisher interface {
	Publish(ctx context.Context, event Event) error
}
