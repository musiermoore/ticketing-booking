package kafka

import (
	"context"

	"github.com/segmentio/kafka-go"
)

func NewProducer(broker, topic string) *kafka.Writer {
	return &kafka.Writer{
		Addr:     kafka.TCP(broker),
		Topic:    topic,
		Balancer: &kafka.LeastBytes{},
	}
}

func PublishBookingEvent(producer *kafka.Writer, bookingID string) error {
	msg := kafka.Message{
		Key:   []byte(bookingID),
		Value: []byte("Booking created: " + bookingID),
	}
	return producer.WriteMessages(context.Background(), msg)
}
