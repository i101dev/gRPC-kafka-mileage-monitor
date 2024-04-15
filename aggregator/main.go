package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"

	"github.com/microservices/types"
)

func main() {

	listenAddr := flag.String("listenaddr", ":3000", "this listen addr of the HTTP server")

	flag.Parse()

	store := NewMemoryStore()

	var (
		svc = NewInvoiceAggregator(store)
	)

	makeHTTPTransport(*listenAddr, svc)

	fmt.Println("invoicer is online and working fine!")
}

func makeHTTPTransport(listenAddr string, svc Aggregator) {
	fmt.Println("HTTP transport running on port", listenAddr)
	http.HandleFunc("/aggregate", handleAggregate(svc))
	http.ListenAndServe(listenAddr, nil)
}

func handleAggregate(svc Aggregator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var distance types.Distance
		if err := json.NewDecoder(r.Body).Decode(&distance); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}
}
