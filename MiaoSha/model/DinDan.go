package model

import (
	"MiaoSha/dao"
	"time"
)

type DinDan struct {
	UserId    string
	Voucher   string
	CreatTime string
}

func CreatDd(userid string, voucher string) {
	// 创建订单
	order := DinDan{
		UserId:    "123",
		Voucher:   "abc123",
		CreatTime: time.Now().Format("2006-01-02 15:04:05"),
	}
	dao.DB.Create(&order)
}
