package main

import (
	"context"
	"log"
	"time"

	"github.com/microservices/aggregator/client"
	"github.com/microservices/types"
)

func main() {
	log.Print("Check 1")

	c, err := client.NewGRPCClient(":3001")

	log.Print("Check 2")

	if err != nil {
		log.Fatal(err)
	}

	log.Print("Check 3")

	if _, err := c.Aggregate(context.TODO(), &types.AggregateRequest{
		ObuID: 1,
		Value: 2.34,
		Unix:  time.Now().UnixNano(),
	}); err != nil {
		log.Fatal("\n*** >>> failed [c.Aggregate] -", err)
	}

	log.Print("Check 4")
}
