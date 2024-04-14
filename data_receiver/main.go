package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/gorilla/websocket"
	"github.com/microservices/types"
)

var kafkaTopic = "obudata"

func main() {
	recv, err := NewDataReceiver()

	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/ws", recv.handleWS)
	http.ListenAndServe(":30000", nil)
}

type DataReceiver struct {
	msgch chan types.OBUData
	conn  *websocket.Conn
	prod  *kafka.Producer
}

func NewDataReceiver() (*DataReceiver, error) {

	p, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": "localhost"})

	if err != nil {
		return nil, err
	}

	// Start another goroutine to check if we have delivered the data.
	go func() {
		for e := range p.Events() {
			switch ev := e.(type) {
			case *kafka.Message:
				if ev.TopicPartition.Error != nil {
					fmt.Printf("Delivery failed: %v\n", ev.TopicPartition)
				} else {
					fmt.Printf("Delivered message to %v\n", ev.TopicPartition)
				}
			}
		}
	}()

	return &DataReceiver{
		msgch: make(chan types.OBUData, 128),
		prod:  p,
	}, nil
}

func (dr *DataReceiver) produceData(data types.OBUData) error {

	b, err := json.Marshal(data)

	if err != nil {
		return err
	}

	dr.prod.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{
			Topic:     &kafkaTopic,
			Partition: kafka.PartitionAny,
		},
		Value: b,
	}, nil)

	fmt.Printf("Received OBU data from [%d] :: <lat - %.2f> <long - %.2f> \n", data.OBUID, data.Lat, data.Long)

	return nil
}

func (dr *DataReceiver) handleWS(w http.ResponseWriter, r *http.Request) {

	u := websocket.Upgrader{
		ReadBufferSize:  1028,
		WriteBufferSize: 1028,
	}

	conn, err := u.Upgrade(w, r, nil)

	if err != nil {
		log.Fatal(err)
	}

	dr.conn = conn

	go dr.wsReceiveLoop()
}
func (dr *DataReceiver) wsReceiveLoop() {

	// fmt.Println("OBU client connected!")

	for {
		var data types.OBUData

		if err := dr.conn.ReadJSON(&data); err != nil {
			log.Println("*** >>> Error reading JSON:", err)
			continue
		}

		if err := dr.produceData(data); err != nil {
			fmt.Println("\n*** >>> [kafka production error] -", err)
		}

		// fmt.Printf("Received OBU data from [%d] :: <lat - %.2f> <long - %.2f> \n", data.OBUID, data.Lat, data.Long)

		// dr.msgch <- data
	}
}
