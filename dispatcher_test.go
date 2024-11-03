package dispatcher_test

import (
	"context"
	"testing"

	"github.com/adnvilla/dispatcher-go"
	"github.com/adnvilla/dispatcher-go/mock"
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

}
