package main

import (
	"context"
	"encoding/json"
	"time"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/microservices/aggregator/client"
	"github.com/microservices/types"
	"github.com/sirupsen/logrus"
)

type DataConsumer interface {
	ConsumeData()
}

type KafkaConsumer struct {
	isRunning   bool
	consumer    *kafka.Consumer
	aggClient   client.Client
	calcService CalculatorServicer
}

func NewKafkaConsumer(topic string, svc CalculatorServicer, ac client.Client) (*KafkaConsumer, error) {

	// fmt.Println("\n*** >>> [topic] -", topic)

	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": "localhost",
		"auto.offset.reset": "earliest",
		"group.id":          "myGroup",
	})

	if err != nil {
		return nil, err
	}

	c.SubscribeTopics([]string{topic}, nil)

	return &KafkaConsumer{
		consumer:    c,
		calcService: svc,
		aggClient:   ac,
	}, nil
}

func (c *KafkaConsumer) Start() {
	// logrus.Info("\n*** >>> kafka transport started")
	c.isRunning = true
	c.readMessageLoop()
}

func (c *KafkaConsumer) readMessageLoop() {

	for c.isRunning {

		msg, err := c.consumer.ReadMessage(-1)

		// logrus.Info("\n*** >>> [readMessageLoop] -", msg)

		if err != nil {
			logrus.Errorf("\n*** >>> kafka consume error -- %s", err)
			continue
		}

		var data types.OBUData

		if err := json.Unmarshal(msg.Value, &data); err != nil {
			logrus.Errorf("\n*** >>> JSON serialization error -- %s", err)
			continue
		}

		distance, err := c.calcService.CalculateDistance(data)

		if err != nil {
			logrus.Errorf("\n*** >>> calculation error -- %s", err)
			continue
		}

		req := &types.AggregateRequest{
			Unix:  time.Now().UnixNano(),
			ObuID: int32(data.OBUID),
			Value: distance,
		}

		if err := c.aggClient.Aggregate(context.Background(), req); err != nil {
			logrus.Error("\n*** >>> aggregate error:", err)
			continue
		}
	}
}
