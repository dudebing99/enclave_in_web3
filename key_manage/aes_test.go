package key_manage

import (
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"testing"
)

func TestAES(t *testing.T) {
	text := "hello from enclave"
	key := []byte("#HvL%$o0oNNoOZnk#o2qbqCeQB1iXeIR")

	fmt.Printf("明文: %s\n秘钥: %s\n", text, string(key))
	encrypted, err := AesEncrypt([]byte(text), key)
	if err != nil {
		panic(err)
	}
	fmt.Printf("加密后: %s\n", base64.StdEncoding.EncodeToString(encrypted))
	//encrypted, _ := base64.StdEncoding.DecodeString("xvhqp8bT0mkEcAsNK+L4fw==")
	origin, err := AesDecrypt(encrypted, key)
	if err != nil {
		panic(err)
	}

	fmt.Printf("解密后明文: %s\n", string(origin))
}

func encryptPrivateKey(privateKey string) (hexedEncrypted string, err error) {
	encryptionKey := []byte("!+kuuFokDD0QhtayezsKHDteEyTm!2$$")

	encrypted, err := AesEncrypt([]byte(privateKey), encryptionKey)
	if err != nil {
		return
	}
	hexedEncrypted = hex.EncodeToString(encrypted)

	fmt.Println("encrypted: ", hexedEncrypted)

	return
}

func decryptPrivateKey(hexedEncrypted string) (decrypted string, err error) {
	encryptionKey := []byte("!+kuuFokDD0QhtayezsKHDteEyTm!2$$")

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

func TestAESEx(t *testing.T) {
	key := "aead75071f4a9437df36d08acdcbb78b8dca55d02f0631f33f72274e9ee45a98"
	fmt.Println("original private key: ", key)

	encodedEncryptedPrivateKey, err := encryptPrivateKey(key)
	if err != nil {
		panic(err)
	}
	fmt.Println("encrypted private key: ", encodedEncryptedPrivateKey)
	privateKey, err := decryptPrivateKey("e9cbda7f70d8399a0f6adb0c9b36130d9d2c6181abab554593b90c4c60b3d2a2ffe5936b83327488f16d28cb3efa3535128884f9381d7e13530df30be46f306b8de96ab429ecac9348d05be1e8df6c84")
	if err != nil {
		panic(err)
	}
	fmt.Println("final private key: ", privateKey)

	if key != privateKey {
		panic("seeds dis-matched")
	}
}
