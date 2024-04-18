package main

import (
	"fmt"
	"log"

	"github.com/microservices/aggregator/client"
)

const (
	kafkaTopic         = "obudata"
	aggregatorEndPoint = "http://127.0.0.1:3000"
)

func main() {

	var (
		err error
		svc CalculatorServicer
	)

	svc = NewCalculatorService()
	svc = NewLogMiddleware(svc)

	httpClient := client.NewHTTPClient(aggregatorEndPoint)
	// grpcClient, err := client.NewGRPCClient(aggregatorEndPoint)

	// if err != nil {
	// 	log.Fatal(err)
	// }

	kafkaConsumer, err := NewKafkaConsumer(kafkaTopic, svc, httpClient)

	if err != nil {
		log.Fatal(`Error chk [0x1] -`, err)
	}

	kafkaConsumer.Start()

	fmt.Println("distance calculator online")
}
