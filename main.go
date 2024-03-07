package spanner101

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"cloud.google.com/go/spanner"
	"github.com/sinmetal/spanner101/internal/trace"
	stores1 "github.com/sinmetal/spanner101/pattern1/stores"
	stores2 "github.com/sinmetal/spanner101/pattern2/stores"
	stores3 "github.com/sinmetal/spanner101/pattern3/stores"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

func RunServer() {
	ctx := context.Background()

	fmt.Println("starting server...")

	spannerClientConfig := spanner.ClientConfig{}
	meterProvider := trace.GetMeterProvider()
	if meterProvider != nil {
		spannerClientConfig.OpenTelemetryMeterProvider = meterProvider
	}

	database1 := os.Getenv("SPANNER_DATABASE1")
	sc1, err := spanner.NewClientWithConfig(ctx, database1, spannerClientConfig)
	if err != nil {
		panic(err)
	}
	ordersStore1, err := stores1.NewOrdersStore(sc1)
	if err != nil {
		panic(err)
	}

	database2 := os.Getenv("SPANNER_DATABASE2")
	sc2, err := spanner.NewClientWithConfig(ctx, database2, spannerClientConfig)
	if err != nil {
		panic(err)
	}
	ordersStore2, err := stores2.NewOrdersStore(sc2)
	if err != nil {
		panic(err)
	}

	database3 := os.Getenv("SPANNER_DATABASE3")
	sc3, err := spanner.NewClientWithConfig(ctx, database3, spannerClientConfig)
	if err != nil {
		panic(err)
	}
	ordersStore3, err := stores3.NewOrdersStore(sc3)
	if err != nil {
		panic(err)
	}

	handlers := &Handlers{
		OrdersStore1: ordersStore1,
		OrdersStore2: ordersStore2,
		OrdersStore3: ordersStore3,
	}

	mux := http.NewServeMux()
	mux.Handle("/insert", otelhttp.NewHandler(http.HandlerFunc(handlers.Insert), "/insert"))
	mux.HandleFunc("/", HelloHandler)

	// TODO Shutdown処理

	// Determine port for HTTP service.
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		fmt.Printf("defaulting to port %s", port)
	}

	// Start HTTP server.
	fmt.Printf("listening on port %s", port)
	if err := http.ListenAndServe(":"+port, mux); err != nil {
		log.Fatal(err)
	}
}
