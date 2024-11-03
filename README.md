# dispatcher-go

[![License](https://img.shields.io/github/license/adnvilla/dispatcher-go)](https://github.com/adnvilla/dispatcher-go/blob/main/LICENSE)
[![GitHub Actions](https://github.com/adnvilla/dispatcher-go/actions/workflows/go.yml/badge.svg)](https://github.com/adnvilla/dispatcher-go/actions/workflows/go.yml)
[![Codecov](https://codecov.io/gh/adnvilla/dispatcher-go/branch/main/graph/badge.svg?token=STRT8T67YP)](https://codecov.io/gh/adnvilla/dispatcher-go)
[![Go Report Card](https://goreportcard.com/badge/github.com/adnvilla/dispatcher-go)](https://goreportcard.com/report/github.com/adnvilla/dispatcher-go)

## Overview

Introducing `dispatcher`, a lightweight and extensible command dispatcher for Golang designed to simplify the handling of requests and responses with support for validation and context-based operations. 

## Features

- **Request and Response Handling**: Define and register handlers for custom request types that can execute logic and return structured responses.
- **Context Support**: All operations are executed with Go's `context.Context`, enabling better control over request lifecycles and cancellation.
- **Validation Integration**: Optionally implement `Validator` for request validation, ensuring that invalid requests are caught before processing.
- **Type Safety**: Utilizes Go generics and reflection to ensure type-safe handler registration and execution.
- **Simple API**: Easy-to-use functions for registering handlers, sending requests, and resetting the dispatcher state.

## API Highlights

- `RegisterHandler`: Register a new handler for a specific request type.
- `Send`: Send a request and receive a response, with automatic validation if the handler implements `Validator`.
- `Reset`: Clear all registered handlers, useful for testing scenarios.

## Usage Example

```go
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

func (h *MyHandler) Validate(ctx context.Context, request MyRequest) error {
    if request.Message == "" {
        return fmt.Errorf("message cannot be empty")
    }
    return nil
}

func main() {
    ctx := context.Background()
    handler := &MyHandler{}
    dispatcher.RegisterHandler(handler)

    response, err := dispatcher.Send[MyRequest, MyResponse](ctx, MyRequest{Message: "Hello, world!"})
    if err != nil {
        fmt.Println("Error:", err)
    } else {
        fmt.Println("Response:", response)
    }

    dispatcher.Reset()
}
```

## Known Issues

- Limited error reporting for handler type mismatches.
- Performance optimizations are planned for future releases.

## Future Enhancements

- Improved logging and error handling.
- Support for middleware to add cross-cutting concerns such as logging and metrics.
- Additional examples and documentation.
