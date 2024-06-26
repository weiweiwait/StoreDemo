package api

import (
	"Go-SecKill/service"
	"github.com/gin-gonic/gin"
	"strconv"
)

//不加锁 出现超卖情况

func WithoutLock(c *gin.Context) {
	gid, _ := strconv.Atoi(c.Query("gid"))
	res := service.WithoutLockSecKill(gid)
	c.JSON(res.Status, res)
}

//加锁(sync包中的Mutex类型的互斥锁),没有问题

func WithLock(c *gin.Context) {
	gid, _ := strconv.Atoi(c.Query("gid"))
	res := service.WithLockSecKill(gid)
	c.JSON(res.Status, res)
}

//  加锁(数据库悲观锁，读限定), 出现超卖

func WithPccRead(c *gin.Context) {
	gid, _ := strconv.Atoi(c.Query("gid"))
	res := service.WithPccReadSecKill(gid)
	c.JSON(res.Status, res)
}

// 加锁(数据库悲观锁，改限定), 正常

func WithPccUpdate(c *gin.Context) {
	gid, _ := strconv.Atoi(c.Query("gid"))
	res := service.WithPccUpdateSecKill(gid)
	c.JSON(res.Status, res)
}

// 加锁(数据库乐观锁，正常)

func WithOcc(c *gin.Context) {
	gid, _ := strconv.Atoi(c.Query("gid"))
	res := service.WithOccSecKill(gid)
	c.JSON(res.Status, res)
}

// 使用管道

func WithChannel(c *gin.Context) {
	gid, _ := strconv.Atoi(c.Query("gid"))
	res := service.WithChannelSecKill(gid)
	c.JSON(res.Status, res)
}
