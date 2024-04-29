package utils

import (
	"github.com/gin-gonic/gin"
)

func ParseRequestId(c *gin.Context) (requestId string) {
	v, _ := c.Get("req_id")
	requestId = v.(string)
	return
}
