package dispatcher_test

import (
	"context"
	"fmt"
	"runtime"
	"testing"

	"github.com/adnvilla/dispatcher-go" // Ajusta el import a tu paquete
)

// Tipos de ejemplo para el benchmark
type BenchmarkRequest struct {
	Data string
}

type BenchmarkResponse struct {
	Success bool
}

// Handler de ejemplo para el benchmark
type BenchmarkHandler struct{}

func (h *BenchmarkHandler) Handle(ctx context.Context, request BenchmarkRequest) (BenchmarkResponse, error) {
	return BenchmarkResponse{Success: true}, nil
}

func (h *BenchmarkHandler) Validate(ctx context.Context, request BenchmarkRequest) error {
	if request.Data == "" {
		return nil
	}
	return nil
}

func BenchmarkDispatcher(b *testing.B) {
	handler := &BenchmarkHandler{}
	dispatcher.RegisterHandler(handler)

	ctx := context.Background()
	request := BenchmarkRequest{Data: "benchmark"}
	defer dispatcher.Reset()

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err := dispatcher.Send[BenchmarkRequest, BenchmarkResponse](ctx, request)
		if err != nil {
			b.Errorf("unexpected error: %v", err)
		}
	}
}

func BenchmarkDispatcherConcurrent(b *testing.B) {
	handler := &BenchmarkHandler{}
	dispatcher.RegisterHandler(handler)
	ctx := context.Background()
	request := BenchmarkRequest{Data: "benchmark"}
	defer dispatcher.Reset()

	b.ReportAllocs()
	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, err := dispatcher.Send[BenchmarkRequest, BenchmarkResponse](ctx, request)
			if err != nil {
				b.Errorf("unexpected error: %v", err)
			}
		}
	})
}

func BenchmarkDispatcherCPUs(b *testing.B) {
	handler := &BenchmarkHandler{}
	dispatcher.RegisterHandler(handler)

	ctx := context.Background()
	request := BenchmarkRequest{Data: "benchmark"}
	defer dispatcher.Reset()

	b.ReportAllocs()
	b.ResetTimer()

	for cpus := 1; cpus <= runtime.NumCPU(); cpus++ {
		b.Run(fmt.Sprintf("CPUs=%d", cpus), func(b *testing.B) {
			runtime.GOMAXPROCS(cpus) // Configura los CPUs
			b.RunParallel(func(pb *testing.PB) {
				for pb.Next() {
					_, err := dispatcher.Send[BenchmarkRequest, BenchmarkResponse](ctx, request)
					if err != nil {
						b.Errorf("unexpected error: %v", err)
					}
				}
			})
		})
	}
}
