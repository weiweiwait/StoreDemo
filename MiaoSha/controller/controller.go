package controller

import "github.com/gin-gonic/gin"

func MiaoSha(c *gin.Context) {
	vocher_id := c.Query("vocherid")
	user_id := c.Query("userid")

}
