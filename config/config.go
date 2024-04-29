package config

import (
	"enclave_in_web3/utils"
	"flag"
	"github.com/spf13/viper"
)

var (
	config string
)

func InitConfig() {
	flag.StringVar(&config, "config", "./conf/application.yml", "")
	flag.Parse()

	viper.SetConfigFile(config)
	viper.SetConfigType("yaml")
	err := viper.ReadInConfig()
	utils.CheckError(err)
}
