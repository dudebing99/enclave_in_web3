package data

import (
	"enclave_in_web3/utils"
	"errors"
	"github.com/golang/glog"
	"github.com/spf13/viper"
	"github.com/syndtr/goleveldb/leveldb"
	"sync"
)

var leveldbMgr *LevelDBMgr

var ErrLevelDBConfig = errors.New("leveldb config error")

var ErrLevelDBUninitialized = errors.New("leveldb uninitialized")

func InitLevelDBMgr() {
	leveldbMgr = newLevelDBMgr(viper.Sub("data.leveldb"))
}

func ReleaseLevelDBMgr() {
	if leveldbMgr != nil {
		leveldbMgr.Close()
		leveldbMgr = nil
	}
}

func GetLevelDB(name string) (*leveldb.DB, error) {
	if leveldbMgr == nil {
		panic(ErrLevelDBUninitialized)
	}

	return leveldbMgr.getLevelDB(name)
}

func MustGetLevelDB(name string) *leveldb.DB {
	if leveldbMgr == nil {
		panic(ErrLevelDBUninitialized)
	}

	return leveldbMgr.mustGetLevelDB(name)
}

func newLevelDBMgr(conf *viper.Viper) *LevelDBMgr {
	dbMgr := &LevelDBMgr{
		dbMap:    make(map[string]*leveldb.DB),
		mutex:    &sync.Mutex{},
		dbConfig: conf,
	}
	return dbMgr
}

type LevelDBMgr struct {
	dbMap    map[string]*leveldb.DB
	mutex    *sync.Mutex
	dbConfig *viper.Viper
}

func (mgr *LevelDBMgr) getLevelDB(name string) (*leveldb.DB, error) {
	config := mgr.dbConfig.Sub(name)
	if config == nil {
		return nil, ErrLevelDBConfig
	}

	mgr.mutex.Lock()
	defer mgr.mutex.Unlock()

	db, ok := mgr.dbMap[name]
	if ok {
		return db, nil
	}

	db, err := initLevelDB(config, name)
	if err != nil {
		return nil, err
	}
	mgr.dbMap[name] = db
	return db, nil
}

func (mgr *LevelDBMgr) mustGetLevelDB(name string) *leveldb.DB {
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

	db, err := initLevelDB(config, name)
	utils.CheckError(err)

	mgr.dbMap[name] = db
	return db
}

func (mgr *LevelDBMgr) Close() {
	mgr.mutex.Lock()
	defer mgr.mutex.Unlock()
	for _, db := range mgr.dbMap {
		db.Close()
	}
	mgr.dbMap = make(map[string]*leveldb.DB)
}

func initLevelDB(config *viper.Viper, name string) (*leveldb.DB, error) {
	path := config.GetString("db")
	levelDB, err := leveldb.OpenFile(path, nil)
	if err != nil {
		return nil, err
	}

	glog.Infof("%s db: path: %s", name, path)
	return levelDB, nil
}
