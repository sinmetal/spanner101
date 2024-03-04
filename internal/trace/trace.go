package trace

import (
	"context"
	"errors"
	"fmt"
	"log"

	"cloud.google.com/go/compute/metadata"
	"cloud.google.com/go/spanner"
	texporter "github.com/GoogleCloudPlatform/opentelemetry-operations-go/exporter/trace"
	"go.opentelemetry.io/contrib/detectors/gcp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.24.0"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/api/option"
)

var tracer trace.Tracer

func init() {
	ctx := context.Background()
	fmt.Println("trace init()")

	if metadata.OnGCE() {
		projectID, err := metadata.ProjectID()
		if err != nil {
			log.Fatalf("required google cloud project id: %v", err)
		}

		spanner.EnableOpenTelemetryMetrics()

		exporter, err := texporter.New(
			texporter.WithProjectID(projectID),
			texporter.WithTraceClientOptions(
				[]option.ClientOption{option.WithTelemetryDisabled()}, // otelのtrace送信そのもののtraceは送らない
			),
		)
		if err != nil {
			log.Fatalf("texporter.New: %v", err)
		}
		res, err := resource.New(ctx,
			// Use the GCP resource detector to detect information about the GCP platform
			resource.WithDetectors(gcp.NewDetector()),
			// Keep the default detectors
			resource.WithTelemetrySDK(),
			// Add your own custom attributes to identify your application
			resource.WithAttributes(
				semconv.ServiceNameKey.String("spanner-hands-on"),
				semconv.ServiceVersion("0.1.0"),
			),
		)
		if errors.Is(err, resource.ErrPartialResource) || errors.Is(err, resource.ErrSchemaURLConflict) {
			log.Println(err)
		} else if err != nil {
			log.Fatalf("resource.New: %v", err)
		}
		tp := sdktrace.NewTracerProvider(
			sdktrace.WithSampler(sdktrace.AlwaysSample()), // 1min間に1requestなので、全部出している
			sdktrace.WithBatcher(exporter),
			sdktrace.WithResource(res),
		)
		// TODO Shutdownはどうやろう？ defer tp.Shutdown(ctx) // flushes any pending spans, and closes connections.
		otel.SetTracerProvider(tp)
		tracer = otel.GetTracerProvider().Tracer("github.com/sinmetal/spanner-hands-on")
	}
	if tracer == nil {
		fmt.Println("set default otel tracer")
		tracer = otel.Tracer("github.com/sinmetal/spanner-hands-on")
	}
}

func StartSpan(ctx context.Context, spanName string, ops ...trace.SpanStartOption) (context.Context, trace.Span) {
	return tracer.Start(ctx, spanName, ops...)
}
