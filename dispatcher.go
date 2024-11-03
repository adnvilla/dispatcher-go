package dispatcher

import (
	"context"
	"fmt"
	"reflect"
)

var handlers = make(map[reflect.Type]interface{})

type Request interface{}
type Response interface{}

type Handler[TRequest Request, TResponse Response] interface {
	Handle(ctx context.Context, request TRequest) (TResponse, error)
}

type Validator[TRequest Request] interface {
	Validate(ctx context.Context, request TRequest) error
}

func RegisterHandler[TRequest Request, TResponse Response](handler Handler[TRequest, TResponse]) {
	request := *new(TRequest)
	requestType := reflect.TypeOf(request)

	_, ok := handlers[requestType]
	if ok {
		panic(fmt.Sprintf("Handler already registered %T", request))
	}

	handlers[requestType] = handler
}

func Send[TRequest Request, TResponse Response](ctx context.Context, request TRequest) (TResponse, error) {
	requestType := reflect.TypeOf(request)
	defaultResponse := *new(TResponse)

	handler, ok := handlers[requestType]
	if !ok {
		return defaultResponse, fmt.Errorf("handler not found for %T", request)
	}

	h, ok := handler.(Handler[TRequest, TResponse])
	if !ok {
		return defaultResponse, fmt.Errorf("invalid handler type for request: %T and response: %T", request, defaultResponse)
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

func Reset() {
	handlers = make(map[reflect.Type]interface{})
}
