package eventbus

import (
	"context"
)

type (
	Events struct {
		events []Event
		index  int
	}

	// EventType 事件类型
	EventType string

	// Event 领域事件
	Event interface {
		Type() EventType
	}

	// EventSubscriber 领域事件订阅方
	EventSubscriber interface {
		Handle(ctx context.Context, event Event)
	}
)

func NewEvents() *Events {
	return &Events{
		events: make([]Event, 0),
		index:  0,
	}
}

func (el *Events) Add(event Event) {
	el.events = append(el.events, event)
}

func (el *Events) next() (Event, bool) {
	if el.index >= len(el.events) {
		return nil, false
	}
	el.index++
	return el.events[el.index-1], true
}

func (el *Events) Dispatch(ctx context.Context) {
	for ev, got := el.next(); got; {
		Publish(ctx, ev)
	}
}

func (el *Events) Range(fn func(ev Event) error) {
	for _, event := range el.events {
		if err := fn(event); err != nil {
			return
		}
	}
}
