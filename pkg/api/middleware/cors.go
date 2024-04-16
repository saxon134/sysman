package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// CrossDomain 跨域处理
func CrossDomain() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method

		c.Header("Access-Control-Allow-Origin", c.Request.Header.Get("Origin"))
		c.Header("Access-Control-Allow-Headers", "Content-Type,X-CSRF-Token,auth,uid,sign,t,timestamp")
		c.Header("Access-Control-Allow-Methods", "POST,GET,OPTIONS")
		c.Header("Access-Control-Expose-Headers", "Content-Length,Access-Control-Allow-Origin,Access-Control-Allow-Headers,Content-Type")
		c.Header("Access-Control-Allow-Credentials", "true")

		if method == "OPTIONS" {
			c.AbortWithStatus(http.StatusOK)
			return
		}
		c.Next()
	}
}
