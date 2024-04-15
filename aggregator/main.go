package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"strconv"

	"github.com/microservices/types"
)

func main() {

	listenAddr := flag.String("listenaddr", ":3000", "this listen addr of the HTTP server")

	flag.Parse()

	var (
		store = NewMemoryStore()
		agg   = NewInvoiceAggregator(store)
	)

	agg = NewLogMiddleware(agg)

	makeHTTPTransport(*listenAddr, agg)

	fmt.Println("invoicer is online and working fine!")
}

func makeHTTPTransport(listenAddr string, agg Aggregator) {
	fmt.Println("HTTP transport running on port", listenAddr)
	http.HandleFunc("/aggregate", handleAggregate(agg))
	http.HandleFunc("/invoice", handleGetInvoice(agg))
	http.ListenAndServe(listenAddr, nil)
}

func handleGetInvoice(agg Aggregator) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		values, ok := r.URL.Query()["obu"]

		if !ok {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "missing <obuID>"})
			return
		}

		obuID, err := strconv.Atoi(values[0])

		if err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid <obuID>"})
			return
		}

		invoice, err := agg.CalculateInvoice(obuID)

		if err != nil {
			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to generate invoice"})
			return
		}

		writeJSON(w, http.StatusOK, invoice)
	}
}

func handleAggregate(agg Aggregator) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		var distance types.Distance

		if err := json.NewDecoder(r.Body).Decode(&distance); err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
			return
		}

		if err := agg.AggregateDistance(distance); err != nil {
			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
			return
		}
	}
}

func writeJSON(w http.ResponseWriter, status int, v any) error {

	w.WriteHeader(status)
	w.Header().Add("Content-Type", "application/json")

	return json.NewEncoder(w).Encode(v)
}
