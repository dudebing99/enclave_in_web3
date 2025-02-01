package middleware

import (
	"enclave_in_web3/utils"
	"github.com/gin-gonic/gin"
)

func GenerateRequestId() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("req_id", utils.GenerateRequestId())

		c.Next()
	}
}
