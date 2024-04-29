package data

import (
	"enclave_in_web3/utils"
	"errors"
	"github.com/golang/glog"
	"github.com/spf13/viper"
	"net"
	"net/http"
	"sync"
	"time"
)

var httpMgr *HttpMgr

var ErrHttpConfig = errors.New("http config error")

var ErrHttpUninitialized = errors.New("http uninitialized")

func InitHttpMgr() {
	httpMgr = newHttpMgr(viper.Sub("data.http"))
}

func ReleaseHttpMgr() {
	if httpMgr != nil {
		httpMgr.Close()
		httpMgr = nil
	}
}

func GetHttp(name string) (*http.Client, error) {
	if httpMgr == nil {
		panic(ErrHttpUninitialized)
	}

	return httpMgr.getHttp(name)
}

func MustGetHttp(name string) *http.Client {
	if httpMgr == nil {
		panic(ErrHttpUninitialized)
	}

	return httpMgr.mustGetHttp(name)
}

func newHttpMgr(conf *viper.Viper) *HttpMgr {
	httpMgr := &HttpMgr{
		httpMap:    make(map[string]*http.Client),
		mutex:      &sync.Mutex{},
		httpConfig: conf,
	}
	return httpMgr
}

type HttpMgr struct {
	httpMap    map[string]*http.Client
	mutex      *sync.Mutex
	httpConfig *viper.Viper
}

func (mgr *HttpMgr) getHttp(name string) (*http.Client, error) {
	config := mgr.httpConfig.Sub(name)
	if config == nil {
		return nil, ErrHttpConfig
	}

	mgr.mutex.Lock()
	defer mgr.mutex.Unlock()

	http, ok := mgr.httpMap[name]
	if ok {
		return http, nil
	}

	http, err := initHttp(config, name)
	if err != nil {
		return nil, err
	}
	mgr.httpMap[name] = http
	return http, nil
}

func (mgr *HttpMgr) mustGetHttp(name string) *http.Client {
	config := mgr.httpConfig.Sub(name)
	if config == nil {
		panic(ErrHttpConfig)
	}

	mgr.mutex.Lock()
	defer mgr.mutex.Unlock()

	http, ok := mgr.httpMap[name]
	if ok {
		return http
	}

	http, err := initHttp(config, name)
	utils.CheckError(err)

	mgr.httpMap[name] = http
	return http
}

func (mgr *HttpMgr) Close() {
	mgr.mutex.Lock()
	defer mgr.mutex.Unlock()
	for _, http := range mgr.httpMap {
		http.CloseIdleConnections()
	}
	mgr.httpMap = make(map[string]*http.Client)
}

func initHttp(config *viper.Viper, name string) (*http.Client, error) {
	http := &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyFromEnvironment,
			DialContext: (&net.Dialer{
				Timeout:   30 * time.Second,
				KeepAlive: 30 * time.Second,
			}).DialContext,
			MaxIdleConns: config.GetInt("max-idle-conn"),
			//MaxIdleConnsPerHost: config.GetInt("max-idle-conn-per-host"),
			IdleConnTimeout:    90 * time.Second,
			DisableCompression: config.GetBool("disable-compression"),
		},
		Timeout: 30 * time.Second,
	}

	glog.Infof("%s http: maxIdleConn:%d, disableCompression: %t",
		name, config.GetInt("max-idle-conn"), config.GetBool("disable-compressed"))

	return http, nil
}
