package model

import (
	"context"
	"fmt"
	"github.com/segmentio/kafka-go"
	"log"
	"net"
	"strconv"
	"time"
)

type Message struct {
	UserID    string
	VarcharID string
}
type Order struct {
	UserID    string
	VarcharID string
}

// 创建对应主题
func CreatTopic() {
	// 指定要创建的topic名称
	topic := "order"

	// 连接至任意kafka节点
	conn, err := kafka.Dial("tcp", "localhost:9092")
	if err != nil {
		panic(err.Error())
	}
	defer conn.Close()

	// 获取当前控制节点信息
	controller, err := conn.Controller()
	if err != nil {
		panic(err.Error())
	}
	var controllerConn *kafka.Conn
	// 连接至leader节点
	controllerConn, err = kafka.Dial("tcp", net.JoinHostPort(controller.Host, strconv.Itoa(controller.Port)))
	if err != nil {
		panic(err.Error())
	}
	defer controllerConn.Close()

	topicConfigs := []kafka.TopicConfig{
		{
			Topic:             topic,
			NumPartitions:     1,
			ReplicationFactor: 1,
		},
	}

	// 创建topic
	err = controllerConn.CreateTopics(topicConfigs...)
	if err != nil {
		panic(err.Error())
	} else {
		fmt.Println("g")
	}
}

// 发送订单
func PubList(userid string, varcharid string) {
	topic := "order"
	partition := 0

	// 连接至Kafka集群的Leader节点
	conn, err := kafka.DialLeader(context.Background(), "tcp", "localhost:9092", topic, partition)
	if err != nil {
		log.Fatal("failed to dial leader:", err)
	}
	message := Message{
		UserID:    userid,
		VarcharID: varcharid,
	}
	// 设置发送消息的超时时间
	conn.SetWriteDeadline(time.Now().Add(10 * time.Second))

	// 发送消息
	_, err = conn.WriteMessages(
		kafka.Message{
			Key:   []byte(fmt.Sprintf("message-%d", message.UserID)),
			Value: []byte(message.VarcharID),
		},
	)
	if err != nil {
		log.Fatal("failed to write messages:", err)
	}

	// 关闭连接
	if err := conn.Close(); err != nil {
		log.Fatal("failed to close writer:", err)
	}
}

// 消费订单
// readByConn 连接至kafka后接收消息
func XiaoFei() {
	// Kafka配置
	config := kafka.ReaderConfig{
		Brokers: []string{"localhost:9092"},
		Topic:   "order",
	}

	// 创建Kafka读取器
	reader := kafka.NewReader(config)

	// 从Kafka读取消息并处理
	for {
		msg, err := reader.ReadMessage(context.Background())
		if err != nil {
			log.Fatal("Failed to read message:", err)
		}

		// 解析消息
		userIDStr := string(msg.Key)
		if err != nil {
			log.Println("Failed to parse user ID:", err)
			continue
		}

		varcharID := string(msg.Value)

		// 构建订单对象
		order := Order{
			UserID:    userIDStr,
			VarcharID: varcharID,
		}

		// 存储订单到数据库
		CreatDd(order.UserID, order.VarcharID)
	}

	// 关闭Kafka读取器
	reader.Close()
}
