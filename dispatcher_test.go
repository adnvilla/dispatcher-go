package dispatcher_test

import (
	"context"
	"errors"
	"sync"
	"testing"

	"github.com/adnvilla/dispatcher-go"
	"github.com/adnvilla/dispatcher-go/mock"
	"github.com/stretchr/testify/assert"
	tmock "github.com/stretchr/testify/mock"
)

type UseCaseTest interface {
	Handle(ctx context.Context, input testInput) (testOutput, error)
}

type testInput struct{}
type testOutput struct{}

func TestDispatcher(t *testing.T) {
	t.Run("Test RegisterRequestHandler", func(t *testing.T) {
		ctx := context.Background()
		dispatcher.ResetRequestHandler()
		dispatcher.RegisterHandler(ctx, mock.NewMockHandler[mock.MockRequest, mock.MockResponse](t))
	})
	t.Run("Test RegisterNotificationHandler", func(t *testing.T) {
		ctx := context.Background()
		dispatcher.ResetNotificationHandler()
		dispatcher.RegisterNotificationHandler(ctx, mock.NewMockNotificationHandler[mock.MockNotification](t))
	})
	t.Run("Test SendRequest", func(t *testing.T) {
		dispatcher.ResetRequestHandler()
		ctx := context.Background()
		input := mock.MockRequest{}
		handler := mock.NewMockHandler[mock.MockRequest, mock.MockResponse](t)

		handler.On("Handle", tmock.Anything, tmock.Anything).Return(mock.MockResponse{}, nil)
		handler.On("Validate", tmock.Anything, tmock.Anything).Return(nil)

		dispatcher.RegisterHandler(ctx, handler)
		_, err := dispatcher.Send[mock.MockRequest, mock.MockResponse](ctx, input)
		if err != nil {
			t.Errorf("Error: %v", err)
		}

		handler.AssertExpectations(t)
	})

	t.Run("Test PublishRequest", func(t *testing.T) {
		dispatcher.ResetNotificationHandler()
		ctx := context.Background()
		input := mock.MockNotification{}
		handler := mock.NewMockNotificationHandler[mock.MockNotification](t)

		handler.On("Handle", tmock.Anything, tmock.Anything).Return(nil)

		dispatcher.RegisterNotificationHandler(ctx, handler)
		err := dispatcher.Publish(ctx, input)
		if err != nil {
			t.Errorf("Error: %v", err)
		}

		handler.AssertExpectations(t)
	})

	t.Run("Test PublishNotification with error", func(t *testing.T) {
		dispatcher.ResetNotificationHandler()
		ctx := context.Background()
		input := mock.MockNotification{}
		handler := mock.NewMockNotificationHandler[mock.MockNotification](t)

		handler.On("Handle", tmock.Anything, tmock.Anything).Return(errors.New("error"))

		dispatcher.RegisterNotificationHandler(ctx, handler)
		err := dispatcher.Publish(ctx, input)
		if err == nil {
			t.Errorf("Error: %v", err)
		}

		handler.AssertExpectations(t)

		assert.EqualError(t, err, "error while handling notification")
	})

	t.Run("Test RegisterRequestHandler with panic", func(t *testing.T) {
		dispatcher.ResetRequestHandler()
		ctx := context.Background()
		dispatcher.RegisterHandler(ctx, mock.NewMockHandler[mock.MockRequest, mock.MockResponse](t))
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("The code did not panic")
			}
		}()
		dispatcher.RegisterHandler(ctx, mock.NewMockHandler[mock.MockRequest, mock.MockResponse](t))
	})

	t.Run("Test RegisterNotificationHandler with panic", func(t *testing.T) {
		dispatcher.ResetNotificationHandler()
		ctx := context.Background()
		dispatcher.RegisterNotificationHandler(ctx, mock.NewMockNotificationHandler[mock.MockNotification](t))

		dispatcher.RegisterNotificationHandler(ctx, mock.NewMockNotificationHandler[mock.MockNotification](t))
	})

	t.Run("Test Request Handler not found", func(t *testing.T) {
		dispatcher.ResetRequestHandler()
		ctx := context.Background()
		input := mock.MockRequest{}
		_, err := dispatcher.Send[mock.MockRequest, mock.MockResponse](ctx, input)
		if err == nil {
			t.Errorf("Error: %v", err)
		}
		assert.EqualError(t, err, "handler not found for mock.MockRequest")
	})

	t.Run("Test Notification Handler not found", func(t *testing.T) {
		dispatcher.ResetNotificationHandler()
		ctx := context.Background()
		input := mock.MockNotification{}
		err := dispatcher.Publish(ctx, input)
		if err == nil {
			t.Errorf("Error: %v", err)
		}
		assert.EqualError(t, err, "handler not found for mock.MockNotification")
	})

	t.Run("Test Invalid Request Handler type", func(t *testing.T) {
		dispatcher.ResetRequestHandler()
		ctx := context.Background()
		input := mock.MockRequest{}
		handler := mock.NewMockHandler[mock.MockRequest, mock.MockResponse](t)
		dispatcher.RegisterHandler(ctx, handler)
		_, err := dispatcher.Send[mock.MockRequest, testOutput](ctx, input)
		if err == nil {
			t.Errorf("Error: %v", err)
		}
		assert.EqualError(t, err, "invalid handler type for request: mock.MockRequest and response: dispatcher_test.testOutput")
	})

	t.Run("Test Validator", func(t *testing.T) {
		dispatcher.ResetRequestHandler()
		ctx := context.Background()
		input := mock.MockRequest{}
		handler := mock.NewMockHandler[mock.MockRequest, mock.MockResponse](t)

		handler.On("Handle", tmock.Anything, tmock.Anything).Return(mock.MockResponse{}, nil)
		handler.On("Validate", tmock.Anything, tmock.Anything).Return(nil)

		dispatcher.RegisterHandler(ctx, handler)

		_, err := dispatcher.Send[mock.MockRequest, mock.MockResponse](ctx, input)
		if err != nil {
			t.Errorf("Error: %v", err)
		}

		handler.AssertExpectations(t)
	})

	t.Run("Test Validator with error", func(t *testing.T) {
		dispatcher.ResetRequestHandler()
		ctx := context.Background()
		input := mock.MockRequest{}
		handler := mock.NewMockHandler[mock.MockRequest, mock.MockResponse](t)

		handler.On("Validate", tmock.Anything, tmock.Anything).Return(errors.New("error"))

		dispatcher.RegisterHandler(ctx, handler)

		_, err := dispatcher.Send[mock.MockRequest, mock.MockResponse](ctx, input)
		if err == nil {
			t.Errorf("Error: %v", err)
		}

		handler.AssertExpectations(t)

		assert.EqualError(t, err, "error")
	})
}

func TestDispatcherRequestConcurrent(t *testing.T) {
	ctx := context.Background()
	handler := &BenchmarkHandler{}
	dispatcher.RegisterHandler(ctx, handler)

	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			ctx := context.Background()
			request := BenchmarkRequest{Data: "test"}
			_, err := dispatcher.Send[BenchmarkRequest, BenchmarkResponse](ctx, request)
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		}()
	}

	wg.Wait()
}

func TestDispatcherNotificationConcurrent(t *testing.T) {
	ctx := context.Background()
	handler := &BenchmarkNotificationHandler{}
	dispatcher.RegisterNotificationHandler(ctx, handler)

	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			ctx := context.Background()
			request := BenchmarkNotification{Data: "test"}
			err := dispatcher.Publish[BenchmarkNotification](ctx, request)
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		}()
	}

	wg.Wait()
}
