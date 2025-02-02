package key_manage

import (
	"enclave_in_web3/utils"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/spf13/viper"
	"sort"
	"strings"
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

	seeds         []string
	encryptionKey []byte
	mtx           sync.Mutex
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

	uuid := utils.GenerateUUID()

	key := &keystore.Key{
		Id:         uuid,
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

func (k *Keeper) IsReady() (err error) {
	count := viper.GetInt("keeper.seed.count")
	if len(k.seeds) != count {
		err = errors.New(fmt.Sprintf("encryption seeds absent, expected %d, got %d", count, len(k.seeds)))
		return
	}

	return nil
}

func (k *Keeper) SetEncryptionSeed(seed string) error {
	k.mtx.Lock()
	defer k.mtx.Unlock()

	count := viper.GetInt("keeper.seed.count")
	minSeedLen := viper.GetInt("keeper.seed.minSeedLen-seed-len")
	maxSeedLen := viper.GetInt("keeper.seed.maxSeedLen-seed-len")
	if len(k.seeds) >= count {
		return errors.New("seeds is exceeded")
	}

	if len(seed) > maxSeedLen || len(seed) < minSeedLen {
		return errors.New(fmt.Sprintf("key's length should between %d and %d", minSeedLen, maxSeedLen))
	}

	k.seeds = append(k.seeds, seed)

	// 生成加密 key
	if len(k.seeds) == count {
		// 对种子排序，设置种子不需要关注顺序
		sort.Slice(k.seeds, func(i, j int) bool {
			return strings.Compare(k.seeds[i], k.seeds[j]) >= 0
		})

		seedsStr := ""
		for i := 0; i < len(k.seeds); i++ {
			seedsStr += k.seeds[i]
		}

		k.encryptionKey = utils.GenerateSHA256(seedsStr)
		//fmt.Println("encryption key: ", hex.EncodeToString(encryptionKey))
	}

	return nil
}

func (k *Keeper) EncryptPrivateKey(privateKey string) (hexedEncrypted string, err error) {
	err = k.IsReady()
	if err != nil {
		return
	}

	encrypted, err := AesEncrypt([]byte(privateKey), k.encryptionKey)
	if err != nil {
		return
	}
	hexedEncrypted = hex.EncodeToString(encrypted)

	//fmt.Println("encrypted: ", hexedEncrypted)

	return
}

func (k *Keeper) DecryptPrivateKey(hexedEncrypted string) (decrypted string, err error) {
	err = k.IsReady()
	if err != nil {
		return
	}

	encrypted, err := hex.DecodeString(hexedEncrypted)
	if err != nil {
		return
	}

	decryptedBytes, err := AesDecrypt(encrypted, k.encryptionKey)
	if err != nil {
		return
	}

	decrypted = string(decryptedBytes)
	//fmt.Println("decrypted: ", decrypted)

	return
}
