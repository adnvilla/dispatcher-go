package main

import (
	"context"
	"fmt"

	dispatcher "github.com/adnvilla/dispatcher-go"
)

type MyRequest struct {
	Message string
}

type MyResponse struct {
	Success bool
}

type MyHandler struct{}

func (h *MyHandler) Handle(ctx context.Context, request MyRequest) (MyResponse, error) {
	return MyResponse{Success: true}, nil
}

func (h *MyHandler) Validate(_ context.Context, request MyRequest) error {
	if request.Message == "" {
		return fmt.Errorf("message cannot be empty")
	}
	return nil
}

func main() {
	ctx := context.Background()
	handler := &MyHandler{}
	// Note: RegisterHandler requires context as the first argument.
	dispatcher.RegisterHandler(ctx, handler)

	response, err := dispatcher.Send[MyRequest, MyResponse](ctx, MyRequest{Message: "Hello, world!"})
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("Response:", response)
	}

	dispatcher.ResetRequestHandler()
}
