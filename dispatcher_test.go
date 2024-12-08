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

type testUseCase struct{}

func TestDispatcher(t *testing.T) {
	t.Run("Test RegisterHandler", func(t *testing.T) {
		dispatcher.Reset()
		dispatcher.RegisterHandler[mock.MockRequest, mock.MockResponse](mock.NewMockHandler[mock.MockRequest, mock.MockResponse](t))
	})
	t.Run("Test Send", func(t *testing.T) {
		dispatcher.Reset()
		ctx := context.Background()
		input := mock.MockRequest{}
		handler := mock.NewMockHandler[mock.MockRequest, mock.MockResponse](t)

		handler.On("Handle", tmock.Anything, tmock.Anything).Return(mock.MockResponse{}, nil)
		handler.On("Validate", tmock.Anything, tmock.Anything).Return(nil)

		dispatcher.RegisterHandler[mock.MockRequest, mock.MockResponse](handler)
		_, err := dispatcher.Send[mock.MockRequest, mock.MockResponse](ctx, input)
		if err != nil {
			t.Errorf("Error: %v", err)
		}

		handler.AssertExpectations(t)
	})
	t.Run("Test RegisterHandler with panic", func(t *testing.T) {
		dispatcher.Reset()
		dispatcher.RegisterHandler[mock.MockRequest, mock.MockResponse](mock.NewMockHandler[mock.MockRequest, mock.MockResponse](t))
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("The code did not panic")
			}
		}()
		dispatcher.RegisterHandler[mock.MockRequest, mock.MockResponse](mock.NewMockHandler[mock.MockRequest, mock.MockResponse](t))
	})

	t.Run("Test Handler not found", func(t *testing.T) {
		dispatcher.Reset()
		ctx := context.Background()
		input := mock.MockRequest{}
		_, err := dispatcher.Send[mock.MockRequest, mock.MockResponse](ctx, input)
		if err == nil {
			t.Errorf("Error: %v", err)
		}
		assert.EqualError(t, err, "handler not found for mock.MockRequest")
	})

	t.Run("Test Invalid Handler type", func(t *testing.T) {
		dispatcher.Reset()
		ctx := context.Background()
		input := mock.MockRequest{}
		handler := mock.NewMockHandler[mock.MockRequest, mock.MockResponse](t)
		dispatcher.RegisterHandler(handler)
		_, err := dispatcher.Send[mock.MockRequest, testOutput](ctx, input)
		if err == nil {
			t.Errorf("Error: %v", err)
		}
		assert.EqualError(t, err, "invalid handler type for request: mock.MockRequest and response: dispatcher_test.testOutput")
	})

	t.Run("Test Validator", func(t *testing.T) {
		dispatcher.Reset()
		ctx := context.Background()
		input := mock.MockRequest{}
		handler := mock.NewMockHandler[mock.MockRequest, mock.MockResponse](t)

		handler.On("Handle", tmock.Anything, tmock.Anything).Return(mock.MockResponse{}, nil)
		handler.On("Validate", tmock.Anything, tmock.Anything).Return(nil)

		dispatcher.RegisterHandler(handler)

		_, err := dispatcher.Send[mock.MockRequest, mock.MockResponse](ctx, input)
		if err != nil {
			t.Errorf("Error: %v", err)
		}

		handler.AssertExpectations(t)
	})

	t.Run("Test Validator with error", func(t *testing.T) {
		dispatcher.Reset()
		ctx := context.Background()
		input := mock.MockRequest{}
		handler := mock.NewMockHandler[mock.MockRequest, mock.MockResponse](t)

		handler.On("Validate", tmock.Anything, tmock.Anything).Return(errors.New("error"))

		dispatcher.RegisterHandler(handler)

		_, err := dispatcher.Send[mock.MockRequest, mock.MockResponse](ctx, input)
		if err == nil {
			t.Errorf("Error: %v", err)
		}

		handler.AssertExpectations(t)

		assert.EqualError(t, err, "error")
	})
}

func TestDispatcherConcurrent(t *testing.T) {
	handler := &BenchmarkHandler{}
	dispatcher.RegisterHandler(handler)

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
