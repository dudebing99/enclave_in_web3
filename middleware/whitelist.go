package middleware

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func IPWhiteList(whitelist []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取请求的 IP 地址
		ip := c.ClientIP()
		// 检查 IP 地址是否在白名单中
		allowed := false
		for _, value := range whitelist {
			if value == ip {
				allowed = true
				break
			}
		}
		// 如果 IP 地址不在白名单中，则返回错误信息
		if !allowed {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": fmt.Sprintf("IP address:%s not allowed", ip)})
			return
		}
		// 允许请求继续访问后续的处理函数
		c.Next()
	}
}
