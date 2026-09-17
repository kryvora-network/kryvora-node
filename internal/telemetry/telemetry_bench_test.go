package telemetry

import (
	"math/rand"
	"net/http"
	"net/http/httptest"
	"testing"
)

func BenchmarkHandleStatusHTTP(b *testing.B) {
	srv := NewServer()
	handler := srv.Handler()
	req := httptest.NewRequest(http.MethodGet, "/status", nil)

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, req)
	}
}

func BenchmarkBackoffJitterComputation(b *testing.B) {
	cfg := DefaultBackoffConfig()
	rng := rand.New(rand.NewSource(42))

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = CalculateBackoff(3, cfg, rng)
	}
}
