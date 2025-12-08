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
	handlers sync.Map
)

func RegisterHandler[TRequest Request, TResponse Response](ctx context.Context, handler Handler[TRequest, TResponse]) {
	request := *new(TRequest)
	requestType := reflect.TypeOf(request)

	_, ok := handlers.LoadOrStore(requestType, handler)
	if ok {
		errMsg := fmt.Sprintf("handler already registered for request: %T", request)
		slog.ErrorContext(ctx, errMsg)
		panic(errors.New(errMsg))
	}
}

func Send[TRequest Request, TResponse Response](ctx context.Context, request TRequest) (TResponse, error) {
	requestType := reflect.TypeOf(request)
	defaultResponse := *new(TResponse)

	handler, ok := handlers.Load(requestType)
	if !ok {
		errMsg := fmt.Sprintf("handler not found for %T", request)
		slog.ErrorContext(ctx, errMsg)
		return defaultResponse, errors.New(errMsg)
	}
	h, ok := handler.(Handler[TRequest, TResponse])
	if !ok {
		errMsg := fmt.Sprintf("invalid handler type for request: %T and response: %T", request, defaultResponse)
		slog.ErrorContext(ctx, errMsg)
		return defaultResponse, errors.New(errMsg)
	}

	if validator, ok := handler.(Validator[TRequest]); ok {
		err := validator.Validate(ctx, request)
		if err != nil {
			return defaultResponse, err
		}
	}

	response, err := h.Handle(ctx, request)
	return response, err
}

func ResetRequestHandler() {
	resetSyncMap(&handlers)
}
