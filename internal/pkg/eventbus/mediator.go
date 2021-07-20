package eventbus

import "context"

type (
	// EventMediator 领域事件调度者
	EventMediator interface {
		Publish(ctx context.Context, event Event)
		Subscribe(EventType, EventSubscriber)
	}

	simpleMediator struct {
		subscribers map[EventType][]EventSubscriber
		concurrency chan struct{}
	}
)

var defaultMediator EventMediator = (*simpleMediator)(nil)

func newSimpleMediator(n int) *simpleMediator {
	return &simpleMediator{
		subscribers: make(map[EventType][]EventSubscriber),
		concurrency: make(chan struct{}, n),
	}
}

func (m *simpleMediator) Publish(ctx context.Context, event Event) {
	if subs, exists := m.subscribers[event.Type()]; exists {
		m.concurrency <- struct{}{}
		go func(ctx context.Context, event Event, subs ...EventSubscriber) {
			defer func() { <-m.concurrency }()

			select {
			case <-ctx.Done():
				// todo log
			default:
				for _, sub := range subs {
					sub.Handle(ctx, event)
				}
			}
		}(ctx, event, subs...)
	}
}

func (m *simpleMediator) Subscribe(t EventType, sub EventSubscriber) {
	if _, exists := m.subscribers[t]; !exists {
		m.subscribers[t] = []EventSubscriber{sub}
		return
	}
	m.subscribers[t] = append(m.subscribers[t], sub)
}

func Publish(ctx context.Context, event Event) {
	defaultMediator.Publish(ctx, event)
}

func Subscribe(t EventType, sub EventSubscriber) {
	defaultMediator.Subscribe(t, sub)
}

func ResetMediator(m EventMediator) {
	defaultMediator = m
}

func init() {
	ResetMediator(newSimpleMediator(10))
}
