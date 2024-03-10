package dao

import (
	"errors"
	"github.com/go-redis/redis"
	"time"
)

func ConnectToRedis() (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379", // Redis 服务器地址
		Password: "",               // Redis 访问密码（如果有的话）
		DB:       0,                // Redis 数据库索引
	})

	// 测试连接是否成功
	_, err := client.Ping().Result()
	if err != nil {
		return nil, err
	}

	return client, nil
}

type RedisLock struct {
	Client  *redis.Client
	LockKey string
	Timeout time.Duration
}

func NewLock(client *redis.Client, lockKey string, timeout time.Duration) *RedisLock {
	return &RedisLock{
		Client:  client,
		LockKey: lockKey,
		Timeout: timeout,
	}
}

func (l *RedisLock) Acquire() error {
	// 尝试获取锁
	success, err := l.Client.SetNX(l.LockKey, "locked", l.Timeout).Result()
	if err != nil {
		return err
	}
	if !success {
		return errors.New("failed to acquire lock")
	}
	return nil
}

func (l *RedisLock) Release() error {
	// 释放锁
	_, err := l.Client.Del(l.LockKey).Result()
	if err != nil {
		return err
	}
	return nil
}
