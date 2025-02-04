package dispatcher

import (
	"context"
	"sync"
)

type Request interface{}
type Response interface{}
type Notification interface{}

type Handler[TRequest Request, TResponse Response] interface {
	Handle(ctx context.Context, request TRequest) (TResponse, error)
}

type NotificationHandler[TNotification Notification] interface {
	Handle(ctx context.Context, notification TNotification) error
}
type Validator[TRequest Request] interface {
	Validate(ctx context.Context, request TRequest) error
}

func resetSyncMap(m *sync.Map) {
	m.Range(func(key, value interface{}) bool {
		m.Delete(key)
		return true
	})
}
