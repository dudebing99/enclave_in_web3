package data

import (
	"enclave_in_web3/utils"
	"errors"
	"fmt"
	_ "github.com/go-sql-driver/mysql" // Do NOT remove me.
	"github.com/golang/glog"
	"github.com/jinzhu/gorm"
	"github.com/spf13/viper"
	"sync"
	"time"
)

var sqlMgr *SQLDBMgr

var ErrSQLConfig = errors.New("sql config error")

var ErrSQLUninitialized = errors.New("sql uninitialized")

func InitSQLMgr() {
	sqlMgr = newSQLDBMgr(viper.Sub("data.mysql"))
}

func ReleaseSQLMgr() {
	if sqlMgr != nil {
		sqlMgr.Close()
		sqlMgr = nil
	}
}

func GetDB(name string) (*gorm.DB, error) {
	if sqlMgr == nil {
		panic(ErrSQLUninitialized)
	}

	return sqlMgr.getDB(name)
}

func MustGetDB(name string) *gorm.DB {
	if sqlMgr == nil {
		panic(ErrSQLUninitialized)
	}

	return sqlMgr.mustGetDB(name)
}

func newSQLDBMgr(conf *viper.Viper) *SQLDBMgr {
	dbMgr := &SQLDBMgr{
		dbMap:    make(map[string]*gorm.DB),
		mutex:    &sync.Mutex{},
		dbConfig: conf,
	}
	return dbMgr
}

type SQLDBMgr struct {
	dbMap    map[string]*gorm.DB
	mutex    *sync.Mutex
	dbConfig *viper.Viper
}

func (mgr *SQLDBMgr) getDB(name string) (*gorm.DB, error) {
	config := mgr.dbConfig.Sub(name)
	if config == nil {
		return nil, ErrSQLConfig
	}

	mgr.mutex.Lock()
	defer mgr.mutex.Unlock()

	db, ok := mgr.dbMap[name]
	if ok {
		return db, nil
	}

	db, err := initDB(config, name)
	if err != nil {
		return nil, err
	}
	mgr.dbMap[name] = db
	return db, nil
}

func (mgr *SQLDBMgr) mustGetDB(name string) *gorm.DB {
	config := mgr.dbConfig.Sub(name)
	if config == nil {
		panic(ErrSQLConfig)
	}

	mgr.mutex.Lock()
	defer mgr.mutex.Unlock()

	db, ok := mgr.dbMap[name]
	if ok {
		return db
	}

	db, err := initDB(config, name)
	utils.CheckError(err)

	mgr.dbMap[name] = db
	return db
}

func (mgr *SQLDBMgr) Close() {
	mgr.mutex.Lock()
	defer mgr.mutex.Unlock()
	for _, db := range mgr.dbMap {
		db.Close()
	}
	mgr.dbMap = make(map[string]*gorm.DB)
}

type DbLogger struct {
}

func (logger *DbLogger) Print(values ...interface{}) {
	//glog.Info(gorm.LogFormatter(values...)...)
	fmt.Println(gorm.LogFormatter(values...)...)
}

func initDB(config *viper.Viper, name string) (*gorm.DB, error) {
	url := config.GetString("url")
	db, err := gorm.Open("mysql", url)
	if err != nil {
		return nil, err
	}

	logger := &DbLogger{}
	db.SetLogger(logger)

	db.LogMode(config.GetBool("log-mode"))
	db.DB().SetConnMaxLifetime(2 * time.Hour)
	maxIdleConn := config.GetInt("max-idle-conn")
	if maxIdleConn != 0 {
		db.DB().SetMaxIdleConns(maxIdleConn)
	}
	maxOpenConn := config.GetInt("max-open-conn")
	if maxOpenConn != 0 {
		db.DB().SetMaxOpenConns(maxOpenConn)
	}

	if err := db.DB().Ping(); err != nil {
		return db, err
	}

	glog.Infof("%s db: maxIdleConn:%d, maxOpenConn: %d", name, maxIdleConn, maxOpenConn)
	return db, nil
}
