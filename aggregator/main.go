package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"strconv"

	"github.com/microservices/types"
	"google.golang.org/grpc"
)

func main() {

	httpAddr := flag.String("httpAddr", ":3000", "this listen addr of the HTTP server")
	grpcAddr := flag.String("grpcAddr", ":3001", "this listen addr of the GRPC server")

	flag.Parse()

	var (
		store = NewMemoryStore()
		agg   = NewInvoiceAggregator(store)
	)

	agg = NewLogMiddleware(agg)

	go func() {
		log.Fatal(makeGRPCTransport(*grpcAddr, agg))
	}()

	log.Fatal(makeHTTPTransport(*httpAddr, agg))
}

func makeGRPCTransport(listenAddr string, svc Aggregator) error {

	fmt.Println("GRPC transport running on port", listenAddr)

	ln, err := net.Listen("tcp", listenAddr)

	if err != nil {
		return err
	}

	defer ln.Close()

	server := grpc.NewServer([]grpc.ServerOption{}...)

	types.RegisterAggregatorServer(server, NewAggregatorGRPCServer(svc))

	return server.Serve(ln)
}

func makeHTTPTransport(listenAddr string, agg Aggregator) error {
	fmt.Println("HTTP transport running on port", listenAddr)
	http.HandleFunc("/aggregate", handleAggregate(agg))
	http.HandleFunc("/invoice", handleGetInvoice(agg))
	return http.ListenAndServe(listenAddr, nil)
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
