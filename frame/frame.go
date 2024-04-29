package frame

import (
	"enclave_in_web3/config"
	"enclave_in_web3/data"
)

// InitFramework 初始化框架
func InitFramework() {
	// 初始化配置文件必须放在最前面，后续初始化可能依赖此
	config.InitConfig()

	//data.InitHttpMgr()
	//data.InitSQLMgr()
	//data.InitRedisMgr()
	data.InitLevelDBMgr()
}

// ReleaseFramework 释放框架
func ReleaseFramework() {
	//data.ReleaseHttpMgr()
	//data.ReleaseSQLMgr()
	//data.ReleaseRedisMgr()
	data.ReleaseLevelDBMgr()
}
