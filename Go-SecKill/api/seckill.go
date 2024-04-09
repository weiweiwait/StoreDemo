package api

import (
	"Go-SecKill/service"
	"github.com/gin-gonic/gin"
	"strconv"
)

func WithoutLock(c *gin.Context) {
	gid, _ := strconv.Atoi(c.Query("gid"))
	res := service.WithoutLockSecKill(gid)
	c.JSON(res.Status, res)
}
