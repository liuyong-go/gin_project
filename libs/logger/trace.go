package logger

import (
	"context"
	"encoding/hex"
	"fmt"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc/metadata"
)

const (
	traceIDStr = "4bf92f3577b34da6a3ce929d0e0e4736" //随机，长度32
	spanIDStr  = "00f067aa0ba902b7"                 //随机长度16
)

type Trace struct {
	TraceID string
	SpanID  string
}

type metadataSupplier struct {
	metadata *metadata.MD
}

func mustTraceIDFromHex(s string) (t trace.TraceID) {
	t, err := trace.TraceIDFromHex(s)
	fmt.Println("TRACE", t.String())
	if err != nil {
		return
	}
	return
}

func mustSpanIDFromHex(s string) (t trace.SpanID) {
	t, err := trace.SpanIDFromHex(s)
	if err != nil {
		return
	}
	return
}

//注入新的trace stateStr : "key1=value1,key2=value2"
func InjectTraceContext(ctx context.Context, preTraceID string, stateStr string) (ltrace *Trace, err error) {
	var state trace.TraceState
	if stateStr != "" {
		state, err = trace.ParseTraceState(stateStr)
		if err != nil {
			return
		}
	}
	var traceID trace.TraceID
	if preTraceID == "" {
		traceID = mustTraceIDFromHex(traceIDStr)
	} else {
		var decoded []byte
		decoded, err = hex.DecodeString(preTraceID)
		if err != nil {
			return
		}
		copy(traceID[:], decoded)
	}
	fmt.Println("traceid", traceID.String())
	sc := trace.NewSpanContext(trace.SpanContextConfig{
		TraceID:    traceID,
		SpanID:     mustSpanIDFromHex(spanIDStr),
		TraceState: state,
		Remote:     true,
	})
	ctx = trace.ContextWithRemoteSpanContext(ctx, sc)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{}, propagation.Baggage{}))
	propagator := otel.GetTextMapPropagator()
	md := metadata.MD{}
	propagator.Inject(ctx, &metadataSupplier{
		metadata: &md,
	})
	spanCtx := trace.SpanContextFromContext(ctx)
	fmt.Println("trace", spanCtx.TraceID().String())
	fmt.Println("span", spanCtx.SpanID().String())
	ltrace = &Trace{
		TraceID: spanCtx.TraceID().String(),
		SpanID:  spanCtx.SpanID().String(),
	}
	return
}
func (s *metadataSupplier) Get(key string) string {
	values := s.metadata.Get(key)
	if len(values) == 0 {
		return ""
	}

	return values[0]
}

func (s *metadataSupplier) Set(key, value string) {
	s.metadata.Set(key, value)
}

func (s *metadataSupplier) Keys() []string {
	out := make([]string, 0, len(*s.metadata))
	for key := range *s.metadata {
		out = append(out, key)
	}

	return out
}
