package router

import (
	api2 "enclave_in_web3/api"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

// Setup sets up all controllers.
func Setup(router *gin.Engine) {

	relativePath := "api"
	enableWhiteList := viper.GetBool("gateway.enable-white-list")
	if !enableWhiteList {
		//echo "enclave"|md5 = 8ed8d7fe15437d09aeba8b757cc14cdc
		//relativePath = "8ed8d7fe15437d09aeba8b757cc14cdc/api"
	}
	api := router.Group(relativePath)

	sign := api2.SignController{Router: api}
	sign.Setup()

	key := api2.KeyController{Router: api}
	key.Setup()

	version := api2.HealthController{Router: api}
	version.Setup()
}
