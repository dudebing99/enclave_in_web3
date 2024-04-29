package api

import (
	"enclave_in_web3/dtos"
	"enclave_in_web3/utils"
	"enclave_in_web3/vsock"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"net/http"
)

// SignController 离线签名
type SignController struct {
	Router *gin.RouterGroup
}

// Setup 初始化路由
func (controller *SignController) Setup() {
	handle := controller.Router.Group("sign")
	{
		handle.POST("message", controller.signMessage)
		handle.POST("transaction", controller.signTransaction)
	}
}

func (controller *SignController) signMessage(c *gin.Context) {
	var req dtos.SignMessageReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"req_id":     utils.ParseRequestId(c),
			"error_code": utils.InvalidParameter,
			"error_msg":  "invalid params",
			"data":       nil})
		return
	}

	reqJson, _ := json.Marshal(req)
	msgType, rspJson, err := vsock.Process(utils.SignMessageReq, reqJson)
	if err != nil || msgType != utils.SignMessageRsp {
		// Enclave Keeper 内部错误特化处理
		if msgType == utils.InternalErrorType {
			var rsp dtos.InternalError
			json.Unmarshal(rspJson, &rsp)

			c.JSON(http.StatusOK, gin.H{
				"req_id":     utils.ParseRequestId(c),
				"error_code": utils.UpstreamError,
				"error_msg":  rsp.ErrorMsg,
				"data":       nil})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"req_id":     utils.ParseRequestId(c),
			"error_code": utils.InternalError,
			"error_msg":  "internal error",
			"data":       nil})
		return
	}

	var rsp dtos.SignMessageRsp
	json.Unmarshal(rspJson, &rsp)
	c.JSON(http.StatusOK, gin.H{
		"req_id":     utils.ParseRequestId(c),
		"error_code": utils.Success,
		"error_msg":  "ok",
		"data":       rsp})
	return
}

func (controller *SignController) signTransaction(c *gin.Context) {
	var req dtos.SignTransactionReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"req_id":     utils.ParseRequestId(c),
			"error_code": utils.InvalidParameter,
			"error_msg":  "invalid params",
			"data":       nil})
		return
	}

	reqJson, _ := json.Marshal(req)
	msgType, rspJson, err := vsock.Process(utils.SignTransactionReq, reqJson)
	if err != nil || msgType != utils.SignTransactionRsp {
		// Enclave Keeper 内部错误特化处理
		if msgType == utils.InternalErrorType {
			var rsp dtos.InternalError
			json.Unmarshal(rspJson, &rsp)

			c.JSON(http.StatusOK, gin.H{
				"req_id":     utils.ParseRequestId(c),
				"error_code": utils.UpstreamError,
				"error_msg":  rsp.ErrorMsg,
				"data":       nil})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"req_id":     utils.ParseRequestId(c),
			"error_code": utils.InternalError,
			"error_msg":  "internal error",
			"data":       nil})
		return
	}

	var rsp dtos.SignTransactionRsp
	json.Unmarshal(rspJson, &rsp)
	c.JSON(http.StatusOK, gin.H{
		"req_id":     utils.ParseRequestId(c),
		"error_code": utils.Success,
		"error_msg":  "ok",
		"data":       rsp})
	return
}
