package main

import (
	"context"
	"fmt"
	dispatcher "github.com/adnvilla/dispatcher-go"
)

type MyNotification struct {
	Message string
}

type MyNotificationHandler struct{}

func (h *MyNotificationHandler) Handle(ctx context.Context, notification MyNotification) error {
	fmt.Println("Notification received:", notification.Message)
	return nil
}

func main() {
	ctx := context.Background()
	handler := &MyNotificationHandler{}
	dispatcher.RegisterNotificationHandler(ctx, handler)

	err := dispatcher.Publish(ctx, MyNotification{Message: "Something happened!"})
	if err != nil {
		fmt.Println("Error:", err)
	}

	dispatcher.ResetNotificationHandler()
}

