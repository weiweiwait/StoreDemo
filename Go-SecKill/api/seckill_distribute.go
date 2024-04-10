package api

import (
	"Go-SecKill/service"
	"github.com/gin-gonic/gin"
	"strconv"
)

// 基于redis的redission分布式,正常

func WithRedission(c *gin.Context) {
	gid, _ := strconv.Atoi(c.Query("gid"))
	res := service.WithRedissionSecKill(gid)
	c.JSON(res.Status, res)
}

// 基于ETCD的锁, 正常

func WithETCD(c *gin.Context) {
	gid, _ := strconv.Atoi(c.Query("gid"))
	res := service.WithETCDSecKill(gid)
	c.JSON(res.Status, res)
}
