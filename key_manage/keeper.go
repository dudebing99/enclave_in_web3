package key_manage

import (
	"enclave_in_web3/utils"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/google/uuid"
	"sync"
)

type EnclaveManagedKey struct {
	KeyId      string
	Address    string
	PrivateKey string
}

type Keeper struct {
	// 私钥保管箱
	keeper sync.Map
	number uint32
}

func NewKeeper() *Keeper {
	return &Keeper{
		number: 0,
	}
}

func (k *Keeper) AddKey(enclaveManagedKey EnclaveManagedKey, ignoreExisted bool) error {
	keyId := enclaveManagedKey.KeyId

	// 设置存储私钥阈值，超过阈值不处理
	if k.number >= utils.MaxPrivateKeys {
		return errors.New("too much private keys")
	}

	// 是否已存在
	if _, ok := k.keeper.Load(keyId); ok {
		// 如果存在，是否接受忽略
		if ignoreExisted {
			return nil
		} else {
			return errors.New("private key referred to the key id exists")
		}
	}

	k.keeper.Store(keyId, enclaveManagedKey)

	k.number++

	return nil
}

func (k *Keeper) Get(keyId string) (privateKey string, err error) {
	if v, ok := k.keeper.Load(keyId); ok {
		return v.(EnclaveManagedKey).PrivateKey, nil
	}

	return "", errors.New(fmt.Sprintf("%s not found", keyId))
}

func (k *Keeper) Generate() EnclaveManagedKey {
	privateKey, err := crypto.GenerateKey()
	utils.CheckError(err)

	// Create the keyfile object with a random UUID.
	UUID, err := uuid.NewRandom()
	utils.CheckError(err)

	key := &keystore.Key{
		Id:         UUID,
		Address:    crypto.PubkeyToAddress(privateKey.PublicKey),
		PrivateKey: privateKey,
	}

	enclaveManagedKey := EnclaveManagedKey{
		KeyId:      key.Id.String(),
		Address:    key.Address.Hex(),
		PrivateKey: hex.EncodeToString(crypto.FromECDSA(key.PrivateKey)),
	}

	return enclaveManagedKey
}
