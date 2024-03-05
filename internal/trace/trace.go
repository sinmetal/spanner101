package trace

import (
	"context"
	"errors"
	"fmt"
	"log"

	"cloud.google.com/go/spanner"
	texporter "github.com/GoogleCloudPlatform/opentelemetry-operations-go/exporter/trace"
	gcppropagator "github.com/GoogleCloudPlatform/opentelemetry-operations-go/propagator"
	metadatabox "github.com/sinmetalcraft/gcpbox/metadata/cloudrun"
	"go.opentelemetry.io/contrib/detectors/gcp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
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

	if metadatabox.OnCloudRun() {
		installPropagators()

		projectID, err := metadatabox.ProjectID()
		if err != nil {
			log.Fatalf("required google cloud project id: %v", err)
		}

		runService, err := metadatabox.Service()
		if err != nil {
			log.Fatalf("required cloud run service: %v", err)
		}

		revision, err := metadatabox.Revision()
		if err != nil {
			log.Fatalf("required cloud run revision: %v", err)
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
				semconv.ServiceNameKey.String(fmt.Sprintf("spanner-hands-on/%s", runService)),
				semconv.ServiceVersion(revision),
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

func installPropagators() {
	otel.SetTextMapPropagator(
		propagation.NewCompositeTextMapPropagator(
			// Putting the CloudTraceOneWayPropagator first means the TraceContext propagator
			// takes precedence if both the traceparent and the XCTC headers exist.
			gcppropagator.CloudTraceOneWayPropagator{},
			propagation.TraceContext{},
			propagation.Baggage{},
		))
}

func StartSpan(ctx context.Context, spanName string, ops ...trace.SpanStartOption) (context.Context, trace.Span) {
	return tracer.Start(ctx, spanName, ops...)
}
