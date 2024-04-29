package key_manage

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"enclave_in_web3/utils"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/spf13/viper"
	"sort"
	"strings"
	"sync"
)

var seeds []string
var encryptionKey []byte
var mtx sync.Mutex

func IsReady() (err error) {
	count := viper.GetInt("keeper.seed.count")
	if len(seeds) != count {
		err = errors.New(fmt.Sprintf("encryption seeds absent, expected %d, got %d", count, len(seeds)))
		return
	}

	return nil
}

func SetEncryptionSeed(seed string) error {
	mtx.Lock()
	defer mtx.Unlock()

	count := viper.GetInt("keeper.seed.count")
	min := viper.GetInt("keeper.seed.min-seed-len")
	max := viper.GetInt("keeper.seed.max-seed-len")
	if len(seeds) >= count {
		return errors.New("seeds is exceeded")
	}

	if len(seed) > max || len(seed) < min {
		return errors.New(fmt.Sprintf("key's length should between %d and %d", min, max))
	}

	seeds = append(seeds, seed)

	// 生成加密 key
	if len(seeds) == count {
		// 对种子排序，设置种子不需要关注顺序
		sort.Slice(seeds, func(i, j int) bool {
			return strings.Compare(seeds[i], seeds[j]) >= 0
		})

		seedsStr := ""
		for i := 0; i < len(seeds); i++ {
			seedsStr += seeds[i]
		}

		encryptionKey = utils.GenerateSHA256(seedsStr)
		//fmt.Println("encryption key: ", hex.EncodeToString(encryptionKey))
	}

	return nil
}

func PKCS7Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padText...)
}

func PKCS7UnPadding(origData []byte) []byte {
	length := len(origData)
	unPadding := int(origData[length-1])
	return origData[:(length - unPadding)]
}

func AesEncrypt(origData, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockSize := block.BlockSize()
	origData = PKCS7Padding(origData, blockSize)
	blockMode := cipher.NewCBCEncrypter(block, key[:blockSize])
	encrypted := make([]byte, len(origData))
	blockMode.CryptBlocks(encrypted, origData)

	return encrypted, nil
}

func AesDecrypt(encrypted, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockSize := block.BlockSize()
	blockMode := cipher.NewCBCDecrypter(block, key[:blockSize])
	origData := make([]byte, len(encrypted))
	blockMode.CryptBlocks(origData, encrypted)
	origData = PKCS7UnPadding(origData)

	return origData, nil
}

func EncryptPrivateKey(privateKey string) (hexedEncrypted string, err error) {
	err = IsReady()
	if err != nil {
		return
	}

	encrypted, err := AesEncrypt([]byte(privateKey), encryptionKey)
	if err != nil {
		return
	}
	hexedEncrypted = hex.EncodeToString(encrypted)

	fmt.Println("encrypted: ", hexedEncrypted)

	return
}

func DecryptPrivateKey(hexedEncrypted string) (decrypted string, err error) {
	err = IsReady()
	if err != nil {
		return
	}

	encrypted, err := hex.DecodeString(hexedEncrypted)
	if err != nil {
		return
	}

	decryptedBytes, err := AesDecrypt(encrypted, encryptionKey)
	if err != nil {
		return
	}

	decrypted = string(decryptedBytes)
	fmt.Println("decrypted: ", decrypted)

	return
}
