package zipkin

import (
	"context"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/zipkin"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.25.0"
	"go.opentelemetry.io/otel/trace"
	"log"
	"os"
)

type ZipkinProvider struct {
	Url    string
	Logger *log.Logger
}

func NewProvider(url string) *ZipkinProvider {
	return &ZipkinProvider{
		Url:    url,
		Logger: log.New(os.Stderr, "zipkin-example", log.Ldate|log.Ltime|log.Llongfile),
	}
}

// InitTracer creates a new trace provider instance and registers it as global trace provider.
func (zp *ZipkinProvider) InitTracer() (func(context.Context) error, trace.TracerProvider, error) {
	// Create Zipkin Exporter and install it as a global tracer.
	//
	// For demoing purposes, always sample. In a production application, you should
	// configure the sampler to a trace.ParentBased(trace.TraceIDRatioBased) set at the desired
	// ratio.
	exporter, err := zipkin.New(
		zp.Url,
		zipkin.WithLogger(zp.Logger),
	)
	if err != nil {
		return nil, nil, err
	}

	batcher := sdktrace.NewBatchSpanProcessor(exporter)

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithSpanProcessor(batcher),
		sdktrace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceName("zipkin-test"),
		)),
	)
	otel.SetTracerProvider(tp)

	return tp.Shutdown, tp, nil
}
