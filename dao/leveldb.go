package dao

import (
	"enclave_in_web3/data"
)

func Set(k, v []byte, target string) error {
	db := data.MustGetLevelDB(target)
	err := db.Put(k, v, nil)
	return err
}
