package main

import (
	"fmt"

	"github.com/microservices/types"
)

type MemoryStore struct{}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{}
}

func (m *MemoryStore) Insert(d types.Distance) error {
	fmt.Println("inserting distance data to storage")
	return nil
}
