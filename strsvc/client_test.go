package strsvc

import (
	"context"
	"testing"

	"github.com/go-kit/kit/log"
	stdopentracing "github.com/opentracing/opentracing-go"
)

func TestNewHTTPClient(t *testing.T) {
	tracer := stdopentracing.GlobalTracer()
	svc, err := NewHTTPClient("localhost:8080", tracer, log.NewNopLogger())
	ctx := context.Background()

	s, err := svc.Uppercase(ctx, "hello, world")
	if err == nil {
		t.Log(s)
	}
	t.Log(s)
}
