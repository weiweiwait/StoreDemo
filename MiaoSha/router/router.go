package service

//
//import (
//	"MiaoSha/dao"
//	"MiaoSha/model"
//	"context"
//	"errors"
//	"github.com/segmentio/kafka-go"
//	"log"
//	"time"
//)
//
//const (
//	kafkaBrokers = "localhost:9092" // Kafka brokers地址
//	topic        = "order_topic"    // 订单主题名称
//)
//
//func MiaoShaService(vocher_id string, user_id string) string {
//	//连接redis
//	client, err := dao.ConnectToRedis()
//	if err != nil {
//		return "连接到Redis失败"
//	}
//	defer client.Close()
//	//查询是否还有库存
//	sum, err := model.GetStoreSumByID(vocher_id)
//	if err != nil {
//		return "查询失败"
//	}
//	if sum <= 0 {
//		return "没有库存"
//	}
//	// 创建锁对象
//	lock := dao.NewLock(client, "lock:order"+user_id, time.Second*10)
//
//	// 获取锁
//	err = lock.Acquire()
//	if err != nil {
//		return "获取锁失败，请重试"
//	}
//
//	defer func() {
//		// 释放锁
//		err := lock.Release()
//		if err != nil {
//			log.Println("释放锁失败:", err)
//		}
//	}()
//
//	// 获取锁成功，执行秒杀操作
//
//	// TODO: 秒杀逻辑
//
//	return "秒杀成功"
//}
//
//func CreateVoucherOrder(voucherID int64) (string, error) {
//	userID := "未见困"
//
//	// 查询订单
//	count, err := queryOrderCount(userID, voucherID)
//	if err != nil {
//		return "", errors.New("查询订单失败")
//	}
//
//	// 判断是否存在
//	if count > 0 {
//		return "", errors.New("用户已经购买过一次！")
//	}
//
//	// 扣减库存
//	success, err := decreaseStock(voucherID)
//	if err != nil {
//		return "", errors.New("库存不足！")
//	}
//
//	if !success {
//		return "", errors.New("库存不足！")
//	}
//
//	// 创建订单
//	orderID, err := createOrder(userID, voucherID)
//	if err != nil {
//		return "", errors.New("创建订单失败")
//	}
//
//	// 将订单信息发送到Kafka消息队列
//	err = produceOrderEvent(orderID)
//	if err != nil {
//		return "", errors.New("发送订单消息失败")
//	}
//
//	return orderID, nil
//}
//
//// 将订单信息发送到Kafka消息队列
//func produceOrderEvent(orderID string) error {
//	// 连接到Kafka
//	w := kafka.NewWriter(kafka.WriterConfig{
//		Brokers: []string{kafkaBrokers},
//		Topic:   topic,
//	})
//
//	// 构造消息
//	msg := kafka.Message{
//		Key:   []byte(orderID),
//		Value: []byte("order_created"),
//	}
//
//	// 发送消息到Kafka
//	err := w.WriteMessages(context.Background(), msg)
//	if err != nil {
//		return err
//	}
//
//	// 关闭Kafka连接
//	err = w.Close()
//	if err != nil {
//		return err
//	}
//
//	return nil
//}
