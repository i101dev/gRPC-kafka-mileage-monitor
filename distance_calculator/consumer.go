package main

import (
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
	aggClient   *client.Client
	calcService CalculatorServicer
}

func NewKafkaConsumer(topic string, svc CalculatorServicer, ac *client.Client) (*KafkaConsumer, error) {

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
	logrus.Info("\n*** >>> kafka transport started")

	c.isRunning = true
	c.readMessageLoop()
}

func (c *KafkaConsumer) readMessageLoop() {

	for c.isRunning {

		msg, err := c.consumer.ReadMessage(-1)

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

		// fmt.Printf("\ndistance -- %.2f", distance)

		req := types.Distance{
			Value: distance,
			Unix:  time.Now().UnixNano(),
			OBUID: data.OBUID,
		}

		if err := c.aggClient.AggregateInvoice(req); err != nil {
			logrus.Error("aggregate error:", err)
			continue
		}
	}
}
