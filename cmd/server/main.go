package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"cloud.google.com/go/compute/metadata"
	"cloud.google.com/go/spanner"
	stores1 "github.com/sinmetal/spanner101/pattern1/stores"
	stores2 "github.com/sinmetal/spanner101/pattern2/stores"
	stores3 "github.com/sinmetal/spanner101/pattern3/stores"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

func main() {
	ctx := context.Background()
	log.Print("starting server...")

	if metadata.OnGCE() {
		spanner.EnableOpenTelemetryMetrics()
	}

	database1 := os.Getenv("SPANNER_DATABASE1")
	sc1, err := spanner.NewClient(ctx, database1)
	if err != nil {
		panic(err)
	}
	ordersStore1, err := stores1.NewOrdersStore(sc1)
	if err != nil {
		panic(err)
	}

	database2 := os.Getenv("SPANNER_DATABASE2")
	sc2, err := spanner.NewClient(ctx, database2)
	if err != nil {
		panic(err)
	}
	ordersStore2, err := stores2.NewOrdersStore(sc2)
	if err != nil {
		panic(err)
	}

	database3 := os.Getenv("SPANNER_DATABASE3")
	sc3, err := spanner.NewClient(ctx, database3)
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
	mux.HandleFunc("/insert", handlers.Insert)
	mux.HandleFunc("/", HelloHandler)

	// TODO Shutdown処理

	// Determine port for HTTP service.
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("defaulting to port %s", port)
	}

	// Start HTTP server.
	log.Printf("listening on port %s", port)
	if err := http.ListenAndServe(":"+port, otelhttp.NewHandler(mux, "server",
		otelhttp.WithMessageEvents(otelhttp.ReadEvents, otelhttp.WriteEvents),
	)); err != nil {
		log.Fatal(err)
	}
}

func HelloHandler(w http.ResponseWriter, r *http.Request) {
	name := os.Getenv("NAME")
	if name == "" {
		name = "World"
	}
	fmt.Fprintf(w, "Hello %s!\n", name)
}
