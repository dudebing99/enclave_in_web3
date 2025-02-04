package tests

import (
	"enclave_in_web3/key_manage"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"testing"
)

func TestAES(t *testing.T) {
	text := "hello from enclave"
	key := []byte("#HvL%$o0oNNoOZnk#o2qbqCeQB1iXeIR")

	t.Logf("明文: %s", text)
	t.Logf("秘钥: %s", string(key))

	encrypted, err := encryptAES([]byte(text), key)
	if err != nil {
		t.Fatalf("Error encrypting text: %v", err)
	}
	t.Logf("加密后: %s", base64.StdEncoding.EncodeToString(encrypted))

	origin, err := decryptAES(encrypted, key)
	if err != nil {
		t.Fatalf("Error decrypting text: %v", err)
	}
	t.Logf("解密后明文: %s", string(origin))

	// Optional: Add test assertions (if this is a real test case)
	if string(origin) != text {
		t.Errorf("Expected decrypted text '%s', but got '%s'", text, string(origin))
	}
}

func encryptAES(plainText, key []byte) ([]byte, error) {
	return key_manage.AesEncrypt(plainText, key)
}

func decryptAES(encrypted, key []byte) ([]byte, error) {
	return key_manage.AesDecrypt(encrypted, key)
}

func encryptPrivateKey(privateKey string) (hexedEncrypted string, err error) {
	encryptionKey := []byte("!+kuuFokDD0QhtayezsKHDteEyTm!2$$")

	encrypted, err := key_manage.AesEncrypt([]byte(privateKey), encryptionKey)
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

	decryptedBytes, err := key_manage.AesDecrypt(encrypted, encryptionKey)
	if err != nil {
		return
	}

	decrypted = string(decryptedBytes)
	fmt.Println("decrypted: ", decrypted)

	return
}

func TestAESEx(t *testing.T) {
	key := "aead75071f4a9437df36d08acdcbb78b8dca55d02f0631f33f72274e9ee45a98"
	t.Logf("Original private key: %s", key)

	encodedEncryptedPrivateKey, err := encryptPrivateKey(key)
	if err != nil {
		t.Fatalf("Failed to encrypt private key: %v", err)
	}
	t.Logf("Encrypted private key: %s", encodedEncryptedPrivateKey)

	privateKey, err := decryptPrivateKey("2b83a80ca653967e683c2b544d144c917a24ee209c682b0654265468980c1ed027d5da973a5a2c3c59b0be4cc8a731887f65aeec79f4c30351eba2b44c0423ee634fb195559033843a4d8d3942b39691")
	if err != nil {
		t.Fatalf("Failed to decrypt private key: %v", err)
	}
	t.Logf("Final private key: %s", privateKey)

	if key != privateKey {
		t.Errorf("Private keys do not match. Expected: %s, Got: %s", key, privateKey)
	}
}
