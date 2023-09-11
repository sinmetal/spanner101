package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"cloud.google.com/go/spanner"
	stores1 "github.com/sinmetal/spanner101/pattern1/stores"
	stores2 "github.com/sinmetal/spanner101/pattern2/stores"
	stores3 "github.com/sinmetal/spanner101/pattern3/stores"
)

func main() {
	ctx := context.Background()
	log.Print("starting server...")

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
	http.HandleFunc("/insert", handlers.Insert)
	http.HandleFunc("/", HelloHandler)

	// TODO Shutdown処理

	// Determine port for HTTP service.
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("defaulting to port %s", port)
	}

	// Start HTTP server.
	log.Printf("listening on port %s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
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
