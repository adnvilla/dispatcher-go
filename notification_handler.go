package dispatcher

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"reflect"
	"sync"
)

var (
	notifications sync.Map
)

func RegisterNotificationHandler[TNotification Notification](ctx context.Context, handler NotificationHandler[TNotification]) {
	if handler == nil {
		errMsg := "handler cannot be nil"
		slog.ErrorContext(ctx, errMsg)
		return
	}
	notification := *new(TNotification)
	notificationType := reflect.TypeOf(notification)

	handlers, ok := notifications.Load(notificationType)
	if ok {
		h, ok := handlers.([]NotificationHandler[TNotification])
		if !ok {
			h = nil
		}
		notifications.Store(notificationType, append(h, handler))
		return
	}

	notifications.Store(notificationType, []NotificationHandler[TNotification]{handler})
}

func Publish[TNotification Notification](ctx context.Context, notification TNotification) error {
	notificationType := reflect.TypeOf(notification)

	handlers, ok := notifications.Load(notificationType)
	if !ok {
		errMsg := fmt.Sprintf("handler not found for %T", notification)
		slog.ErrorContext(ctx, errMsg)
		return errors.New(errMsg)
	}
	h, ok := handlers.([]NotificationHandler[TNotification])
	if !ok {
		errMsg := fmt.Sprintf("invalid handler type for notification: %T", notification)
		slog.ErrorContext(ctx, errMsg)
		return errors.New(errMsg)
	}

	for _, handler := range h {
		err := handler.Handle(ctx, notification)
		if err != nil {
			errMsg := "error while handling notification"
			return errors.New(errMsg)
		}
	}
	return nil
}

func ResetNotificationHandler() {
	resetSyncMap(&notifications)
}
