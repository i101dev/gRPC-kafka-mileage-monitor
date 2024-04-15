package main

import (
	"fmt"

	"github.com/microservices/types"
)

type Aggregator interface {
	AggregateDistance(types.Distance) error
}

type Storer interface {
	Insert(types.Distance) error
}

type InvoiceAggregator struct {
	store Storer
}

func NewInvoiceAggregator(s Storer) *InvoiceAggregator {
	return &InvoiceAggregator{
		store: s,
	}
}

func (i *InvoiceAggregator) AggregateDistance(distance types.Distance) error {
	fmt.Println("processing distance data")
	return i.store.Insert(distance)
}
