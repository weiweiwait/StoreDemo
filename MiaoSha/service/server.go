package service

import (
	"MiaoSha/dao"
	"MiaoSha/model"
	"context"
	"errors"
	"fmt"
	"github.com/segmentio/kafka-go"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

func handlePendingList() {

	topic := "order"
	partition := 0
	groupID := "group1"

	// 创建Kafka消费者
	config := kafka.ReaderConfig{
		Brokers:   []string{"localhost:9092"},
		GroupID:   groupID,
		Topic:     topic,
		Partition: partition,
	}
	reader := kafka.NewReader(config)
	defer reader.Close()

	// 创建一个等待组，用于等待消费者协程结束
	var wg sync.WaitGroup
	wg.Add(1)

	// 启动消费者协程
	go func() {
		defer wg.Done()

		for {
			// 从主题中读取消息
			msg, err := reader.ReadMessage(context.Background())
			if err != nil {
				fmt.Println("failed to read message:", err)
				break
			}

			// 处理订单信息
			userID := string(msg.Key)
			varcharID := string(msg.Value)
			handleOrder(userID, varcharID)
		}
	}()

	// 等待退出信号
	exitSignal := make(chan os.Signal, 1)
	signal.Notify(exitSignal, syscall.SIGINT, syscall.SIGTERM)
	<-exitSignal

	// 等待消费者协程结束
	wg.Wait()
}
func handleOrder(userId string, varcharId string) {
	client, _ := dao.ConnectToRedis()
	// 存储键值对到Redis
	exist, _ := client.Exists(userId).Result()
	if exist == 0 {
		client.Set(userId, varcharId, time.Second*1000).Err()
		//保存在mysql中
	}

}
func MiaoShaService(vocher_id string, user_id string) string {
	//连接redis
	client, err := dao.ConnectToRedis()
	if err != nil {
		return "连接到Redis失败"
	}
	defer client.Close()
	//查询是否还有库存
	sum, err := model.GetStoreSumByID(vocher_id)
	if err != nil {
		return "查询失败"
	}
	if sum <= 0 {
		return "没有库存"
	}
	// 创建锁对象
	lock := dao.NewLock(client, "lock:order"+user_id, time.Second*10)

	// 获取锁
	err = lock.Acquire()
	if err != nil {
		return "获取锁失败，请重试"
	}

	defer func() {
		// 释放锁
		err := lock.Release()
		if err != nil {
			log.Println("释放锁失败:", err)
		}
	}()

	// 获取锁成功，执行秒杀操作

	// TODO: 秒杀逻辑
	err = model.DecreaseStoreByID(vocher_id)
	if err != nil {
		return "秒杀失败"
	}
	return "秒杀成功"
}

func CreateVoucherOrder(voucherID int64) (string, error) {
	userID := "未见困"

	// 查询订单
	count, err := model.queryOrderCount(userID)
	if err != nil {
		return "", errors.New("查询订单失败")
	}

	// 判断是否存在
	if count > 0 {
		return "", errors.New("用户已经购买过一次！")
	}

	// 扣减库存
	success, err := decreaseStock(voucherID)
	if err != nil {
		return "", errors.New("库存不足！")
	}

	if !success {
		return "", errors.New("库存不足！")
	}

	// 创建订单
	orderID, err := createOrder(userID, voucherID)
	if err != nil {
		return "", errors.New("创建订单失败")
	}

	// 将订单信息发送到Kafka消息队列
	err = produceOrderEvent(orderID)
	if err != nil {
		return "", errors.New("发送订单消息失败")
	}

	return orderID, nil
}
