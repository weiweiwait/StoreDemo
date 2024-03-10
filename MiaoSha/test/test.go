package tset

import (
	"context"
	"fmt"
	"github.com/segmentio/kafka-go"
	"log"
	"testing"
)

func TestsKafka(t *testing.T) {
	// 创建一个 Kafka writer，指定要连接的 brokers 和 topic
	w := kafka.NewWriter(kafka.WriterConfig{
		Brokers:  []string{"localhost:9092"},
		Topic:    "order",
		Balancer: &kafka.LeastBytes{},
	})

	// 创建一个 Kafka message，包含订单的信息
	msg := kafka.Message{
		Key:   []byte("order_id"),
		Value: []byte("123456"),
	}

	// 使用 writer 发送 message
	err := w.WriteMessages(context.Background(), msg)
	if err != nil {
		log.Fatal("failed to write message:", err)
	}

	fmt.Println("message sent successfully")

	// 关闭 writer
	w.Close()
}
