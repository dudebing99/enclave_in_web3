package api

import (
	"enclave_in_web3/utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

// HealthController 版本信息
type HealthController struct {
	Router *gin.RouterGroup
}

// Setup 初始化路由
func (controller *HealthController) Setup() {
	handle := controller.Router.Group("health")
	{
		handle.GET("", controller.health)
		handle.POST("", controller.health)
	}
}

func (controller *HealthController) health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"req_id":     utils.ParseRequestId(c),
		"error_code": utils.Success,
		"error_msg":  "ok",
		"data":       nil,
	})
}
