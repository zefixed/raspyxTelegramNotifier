package kafka

import (
	"context"
	"fmt"
	"github.com/segmentio/kafka-go"
)

type ConsumerConfig struct {
	Topic string
	URL   string
	Group string
}

type Event struct {
	Timestamp string      `json:"timestamp"`
	Message   interface{} `json:"message"`
}

type Consumer struct {
	reader *kafka.Reader
}

func NewConsumer(cfg *ConsumerConfig) *Consumer {
	return &Consumer{
		reader: kafka.NewReader(kafka.ReaderConfig{
			Brokers: []string{cfg.URL},
			Topic:   cfg.Topic,
			GroupID: cfg.Group,
		}),
	}
}

func (c *Consumer) Consume(ctx context.Context, msgChan chan<- kafka.Message) error {
	const op = "kafka.consumer.Consume"

	for {
		msg, err := c.reader.ReadMessage(ctx)
		if err != nil {
			return err
		}
		select {
		case <-ctx.Done():
			return fmt.Errorf("%v: context canceled", op)
		case msgChan <- msg:
		}
	}
}

func (c *Consumer) Close() error {
	return c.reader.Close()
}
