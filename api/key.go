package api

import (
	"enclave_in_web3/dao"
	"enclave_in_web3/dtos"
	"enclave_in_web3/utils"
	"enclave_in_web3/vsock"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"net/http"
)

// KeyController 私钥管理
type KeyController struct {
	Router *gin.RouterGroup
}

// Setup 初始化路由
func (controller *KeyController) Setup() {
	handle := controller.Router.Group("key")
	{
		handle.POST("generate", controller.generateKey)
		handle.POST("add", controller.addKey)

		// 托管钱包相关接口: 设置用户私钥的加密密钥
		handle.POST("set_encryption_seed", controller.setEncryptionSeed)
		// 托管钱包相关接口: 生成用户地址，数据库记录加密后的用户私钥，上层业务服务只需要关联 UID/Key Id
		// TODO: 目前暂只支持以太坊系列地址
		handle.POST("generate_address", controller.generateAddress)
	}
}

func (controller *KeyController) generateKey(c *gin.Context) {
	var req dtos.GenerateKeyReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"req_id":     utils.ParseRequestId(c),
			"error_code": utils.InvalidParameter,
			"error_msg":  "invalid params",
			"data":       nil})
		return
	}

	debug := viper.GetBool("gateway.debug")
	if !debug {
		// 生产模式，不允许返回明文私钥
		showPrivateKey := req.ShowPrivateKey
		if showPrivateKey {
			c.JSON(http.StatusOK, gin.H{
				"req_id":     utils.ParseRequestId(c),
				"error_code": utils.InvalidParameter,
				"error_msg":  "invalid params: show_private_key should be FALSE in production env",
				"data":       nil})
			return
		}
	}
	reqJson, _ := json.Marshal(req)
	msgType, rspJson, err := vsock.Process(utils.GenerateKeyReq, reqJson)
	if err != nil || msgType != utils.GenerateKeyRsp {
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

	var rsp dtos.GenerateKeyRsp
	json.Unmarshal(rspJson, &rsp)
	c.JSON(http.StatusOK, gin.H{
		"req_id":     utils.ParseRequestId(c),
		"error_code": utils.Success,
		"error_msg":  "ok",
		"data":       rsp})
	return
}

func (controller *KeyController) addKey(c *gin.Context) {
	var req dtos.AddKeyReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"req_id":     utils.ParseRequestId(c),
			"error_code": utils.InvalidParameter,
			"error_msg":  "invalid params",
			"data":       nil})
		return
	}

	// 检查参数
	keyId := req.KeyId
	address := req.Address
	privateKey := req.PrivateKey
	if len(keyId) != utils.KeyIdLength ||
		len(address) != utils.AddressLength ||
		len(privateKey) != utils.PrivateKeyLength {
		c.JSON(http.StatusOK, gin.H{
			"req_id":     utils.ParseRequestId(c),
			"error_code": utils.InvalidParameter,
			"error_msg":  "invalid params",
			"data":       nil})
		return
	}

	reqJson, _ := json.Marshal(req)
	msgType, rspJson, err := vsock.Process(utils.AddKeyReq, reqJson)
	if err != nil || msgType != utils.AddKeyRsp {
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

	var rsp dtos.AddKeyRsp
	json.Unmarshal(rspJson, &rsp)
	c.JSON(http.StatusOK, gin.H{
		"req_id":     utils.ParseRequestId(c),
		"error_code": utils.Success,
		"error_msg":  "ok",
		"data":       rsp})
	return
}

func (controller *KeyController) setEncryptionSeed(c *gin.Context) {
	var req dtos.SetEncryptionSeedReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"req_id":     utils.ParseRequestId(c),
			"error_code": utils.InvalidParameter,
			"error_msg":  "invalid params",
			"data":       nil})
		return
	}

	reqJson, _ := json.Marshal(req)
	msgType, rspJson, err := vsock.Process(utils.SetEncryptionSeedReq, reqJson)
	if err != nil || msgType != utils.SetEncryptionSeedRsp {
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

	var rsp dtos.SetEncryptionSeedRsp
	json.Unmarshal(rspJson, &rsp)
	c.JSON(http.StatusOK, gin.H{
		"req_id":     utils.ParseRequestId(c),
		"error_code": utils.Success,
		"error_msg":  "ok",
		"data":       rsp})
	return
}

func (controller *KeyController) generateAddress(c *gin.Context) {
	var req dtos.GenerateAddressReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"req_id":     utils.ParseRequestId(c),
			"error_code": utils.InvalidParameter,
			"error_msg":  "invalid params",
			"data":       nil})
		return
	}

	showPrivateKey := req.ShowPrivateKey
	reqJson, _ := json.Marshal(req)
	msgType, rspJson, err := vsock.Process(utils.GenerateAddressReq, reqJson)
	if err != nil || msgType != utils.GenerateAddressRsp {
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

	// 创建地址成功，保存到本地数据库
	var rsp dtos.GenerateAddressRsp
	json.Unmarshal(rspJson, &rsp)

	k := []byte(rsp.KeyId)
	v := rspJson

	err = dao.Set(k, v, "keystore")
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"req_id":     utils.ParseRequestId(c),
			"error_code": utils.InternalError,
			"error_msg":  "internal error: persist in leveldb",
			"data":       nil})
		return
	}

	// 默认不返回加密之后的私钥给应用层
	if !showPrivateKey {
		rsp.PrivateKey = ""
	}
	c.JSON(http.StatusOK, gin.H{
		"req_id":     utils.ParseRequestId(c),
		"error_code": utils.Success,
		"error_msg":  "ok",
		"data":       rsp})
	return
}
