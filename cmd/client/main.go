package main

import (
	"enclave_in_web3/data"
	"enclave_in_web3/dtos"
	"enclave_in_web3/frame"
	"enclave_in_web3/middleware"
	"enclave_in_web3/router"
	"enclave_in_web3/utils"
	"enclave_in_web3/vsock"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang/glog"
	"github.com/spf13/viper"
	"os"
)

func main() {
	defer glog.Flush()

	fmt.Println("Starting gateway ...")

	frame.InitFramework()
	defer frame.ReleaseFramework()

	// 触发 vsock
	vsock.Trigger()

	if err := loadKeystore(); err != nil {
		utils.CheckError(err)
	}

	// 初始化 web 服务器
	debug := viper.GetBool("gateway.debug")
	if !debug {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()
	// Global middleware
	// Logger middleware will write the logs to gin.DefaultWriter even if you set with GIN_MODE=release.
	// By default, gin.DefaultWriter = os.Stdout
	r.Use(gin.Logger())
	// Recovery middleware recovers from any panics and writes a 500 if there was one.
	r.Use(gin.Recovery())
	//r.Use(static.Serve("/", static.LocalFile("./bin/dist", true)))
	r.Use(middleware.Cors())
	r.Use(middleware.LoggerToFile())

	r.Use(middleware.GenerateRequestId())

	enableWhiteList := viper.GetBool("gateway.enable-white-list")
	if enableWhiteList {
		whiteList := viper.GetStringSlice("gateway.white-list")
		r.Use(middleware.IPWhiteList(whiteList))

		fmt.Printf("white list enabled, in detail: %v\n", whiteList)
	} else {
		fmt.Printf("white list NOT enabled!\n")
	}

	router.Setup(r)

	addr := viper.GetString("gateway.addr")
	var err error
	if len(addr) != 0 {
		err = r.Run(addr)
	} else {
		err = r.Run()
	}
	if err != nil {
		utils.CheckError(err)
	}

	os.Exit(0)
}

func loadKeystore() (err error) {
	loadOnStart := viper.GetBool("persistence-rule.load-on-start")
	if !loadOnStart {
		return
	}

	// 从持久化存储加载数据
	target := "keystore"
	db := data.MustGetLevelDB(target)
	iter := db.NewIterator(nil, nil)
	count := uint32(0)
	debug := viper.GetBool("gateway.debug")
	for iter.Next() {
		k := iter.Key()
		v := iter.Value()
		count += 1

		if debug {
			glog.Infof("loading data, No.%d : %s -> %s\n", count, string(k), string(v))
		}

		reqJson := v
		msgType, rspJson, err := vsock.Process(utils.AddAddressReq, reqJson)
		if err != nil || msgType != utils.AddAddressRsp {
			// Enclave Keeper 内部错误特化处理
			if msgType == utils.InternalErrorType {
				var rsp dtos.InternalError
				json.Unmarshal(rspJson, &rsp)

				return errors.New(fmt.Sprintf("key id: %s, %s", string(k), rsp.ErrorMsg))
			}

			return errors.New(fmt.Sprintf("key id: %s, internal error", string(k)))
		}

		// 无需要处理响应
		//var rsp dtos.AddAddressRsp
		//json.Unmarshal(rspJson, &rsp)
	}

	return
}
