package main

import (
	"fmt"
	"log"
	"time"

	"github.com/go-redis/redis"
)

const (
	LockKey      = "mylock"
	LockValue    = "1"
	LockTimeout  = time.Second * 10
	WatchdogTick = time.Second * 2
)

func main() {
	// 创建Redis客户端
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379", // Redis服务器地址
		Password: "",               // Redis服务器密码
		DB:       0,                // 使用默认数据库
	})
	// 获取锁
	acquired, err := acquireLock(client, LockKey, LockValue, LockTimeout)
	if err != nil {
		log.Fatal("Failed to acquire lock:", err)
	}
	if !acquired {
		log.Println("Failed to acquire lock, another process holds the lock")
		return
	}

	defer func() {
		// 释放锁
		err := releaseLock(client, LockKey, LockValue)
		if err != nil {
			log.Println("Failed to release lock:", err)
		}
	}()

	// 启动看门狗
	go watchdog(client, LockKey, LockValue, LockTimeout, WatchdogTick)

	// 执行需要加锁的逻辑
	fmt.Println("Lock acquired, performing critical section")
	// 在这里执行需要加锁的操作

	// 模拟操作的耗时
	time.Sleep(time.Second * 5)
	fmt.Println("Critical section completed")
}

// 获取锁
func acquireLock(client *redis.Client, lockKey, lockValue string, timeout time.Duration) (bool, error) {
	result, err := client.SetNX(lockKey, lockValue, timeout).Result()
	if err != nil {
		return false, err
	}
	return result, nil
}

// 释放锁
func releaseLock(client *redis.Client, lockKey, lockValue string) error {
	script := `
		if redis.call("GET", KEYS[1]) == ARGV[1] then
			return redis.call("DEL", KEYS[1])
		else
			return 0
		end
	`
	_, err := client.Eval(script, []string{lockKey}, lockValue).Result()
	if err != nil {
		return err
	}
	return nil
}

// 看门狗，定时刷新锁的过期时间
func watchdog(client *redis.Client, lockKey, lockValue string, lockTimeout, tick time.Duration) {
	ticker := time.NewTicker(tick)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			err := client.Expire(lockKey, lockTimeout).Err()
			if err != nil {
				log.Println("Failed to refresh lock expiration:", err)
			}
		}
	}
}
