package mygin

import (
	"github.com/gin-gonic/gin"
	"testing"
)

func Test_Server(t *testing.T){
	r := gin.Default()
	r.GET("/ping", (&PingController{}).Ping)
	r.Run(":8081")
}

type PingController struct {
	*BaseGinController
}

func(this *PingController)Ping(c *gin.Context) {
	c.String(200, "pong")
}
