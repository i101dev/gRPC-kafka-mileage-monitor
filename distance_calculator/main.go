package main

import (
	"fmt"
	"log"

	"github.com/microservices/aggregator/client"
)

const (
	kafkaTopic         = "obudata"
	aggregatorEndPoint = "http://127.0.0.1:3000/aggregate"
)

func main() {

	var (
		err error
		svc CalculatorServicer
	)

	svc = NewCalculatorService()
	svc = NewLogMiddleware(svc)

	client := client.NewHTTPClient(aggregatorEndPoint)

	kafkaConsumer, err := NewKafkaConsumer(kafkaTopic, svc, client)

	if err != nil {
		log.Fatal(err)
	}

	kafkaConsumer.Start()

	fmt.Println("distance calculator online")
}
