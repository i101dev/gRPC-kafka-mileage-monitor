package main

import (
	"fmt"

	"github.com/microservices/types"
)

const basePrice = 3.15

type Aggregator interface {
	AggregateDistance(types.Distance) error
	CalculateInvoice(int) (*types.Invoice, error)
}

type Storer interface {
	Insert(types.Distance) error
	Get(int) (float64, error)
}

type InvoiceAggregator struct {
	store Storer
}

func NewInvoiceAggregator(s Storer) Aggregator {
	return &InvoiceAggregator{
		store: s,
	}
}

func (i *InvoiceAggregator) AggregateDistance(distance types.Distance) error {
	return i.store.Insert(distance)
}
func (i *InvoiceAggregator) CalculateInvoice(obuID int) (*types.Invoice, error) {

	dist, err := i.store.Get(obuID)

	if err != nil {
		return nil, fmt.Errorf("obuID  (%d) - not found", obuID)
	}

	invoice := &types.Invoice{
		OBUID:         obuID,
		TotalDistance: dist,
		TotalAmount:   basePrice * dist,
	}

	return invoice, nil
}
