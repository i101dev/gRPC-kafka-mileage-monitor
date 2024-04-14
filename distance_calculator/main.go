package main

import (
	"fmt"
	"log"
)

func main() {

	var (
		err error
		svc CalculatorServicer
	)

	svc = NewCalculatorService()

	kafkaConsumer, err := NewKafkaConsumer("obudata", svc)

	if err != nil {
		log.Fatal(err)
	}

	kafkaConsumer.Start()

	fmt.Println("distance calculator online")
}
