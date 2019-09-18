package mygin

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"time"

	"github.com/tiptok/gocomm/common"
)

type BaseGinController struct {

}

func(this *BaseGinController)JWTMiddleware()gin.HandlerFunc{
	return func(c *gin.Context){
		token := c.Query("token")
		code := http.StatusOK
		if token == "" {
			code = http.StatusUnauthorized
		} else {
			claims, err := common.ParseJWTToken(token)
			if err != nil {
				code = http.StatusUnauthorized
			} else if time.Now().Unix() > claims.ExpiresAt {
				code = http.StatusUnauthorized
			}
		}
		if code != http.StatusOK {
			this.ResponseJson(c,NewMessage(1).SetHttpCode(code))
			return
		}
		c.Next()
	}
}

//group.Use(Prepare)
func(this *BaseGinController)Prepare(c *gin.Context){
	this.Secure(c)
	this.NoCache(c)
}

// NoCache is a middleware function that appends headers
// to prevent the client from caching the HTTP response.
func (this *BaseGinController)NoCache(c *gin.Context) {
	c.Header("Cache-Control", "no-cache, no-store, max-age=0, must-revalidate, value")
	c.Header("Expires", "Thu, 01 Jan 1970 00:00:00 GMT")
	c.Header("Last-Modified", time.Now().UTC().Format(http.TimeFormat))
	c.Next()
}

// Secure is a middleware function that appends security
// and resource access headers.
func (this *BaseGinController)Secure(c *gin.Context) {
	c.Header("Access-Control-Allow-Origin", "*")
	c.Header("X-Frame-Options", "DENY")
	c.Header("X-Content-Type-Options", "nosniff")
	c.Header("X-XSS-Protection", "1; mode=block")
	c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	c.Header("Access-Control-Allow-Headers", "uid, token,jwt, deviceid, appid,Content-Type,Authorization,from")
	if c.Request.TLS != nil {
		c.Header("Strict-Transport-Security", "max-age=31536000")
	}

	// Also consider adding Content-Security-Policy headers
	// c.Header("Content-Security-Policy", "script-src 'self' https://cdnjs.cloudflare.com")
}

func(this *BaseGinController)GetRequestHead(c *gin.Context)*RequestHead{
	requestHead := &RequestHead{}
	requestHead.Token = c.Query("token")
	requestHead.Version = c.Query("version")
	requestHead.Os = c.Query("os")
	requestHead.From = c.Query("from")
	requestHead.Screen = c.Query("screen")
	requestHead.Model = c.Query("model")
	requestHead.Channel = c.Query("channel")
	requestHead.Net = c.Query("net")
	requestHead.DeviceId = c.Query("deviceid")
	requestHead.Uid, _ = strconv.ParseInt(c.Query("uid"), 10, 64)
	requestHead.AppId, _ = strconv.Atoi(c.Query("appid"))
	requestHead.LoginIp = c.ClientIP()
	requestHead.Jwt = c.Query("jwt")
	return requestHead
}


func(this *BaseGinController)ResponseJson(c *gin.Context,rsp *Message){
	c.JSON(rsp.HttpCode,rsp)
	c.Abort()
}
