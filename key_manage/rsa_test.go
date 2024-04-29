package key_manage

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/hex"
	"encoding/pem"
	"fmt"
	"strings"
	"testing"
)

func TestRSA(t *testing.T) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		panic(err)
	}
	//生成公钥
	publicKey := privateKey.PublicKey

	publicKeyBytes, err := x509.MarshalPKIXPublicKey(&privateKey.PublicKey)
	if err != nil {
		panic(err)
	}
	publicKeyPem := pem.EncodeToMemory(&pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: publicKeyBytes,
	})
	//fmt.Println(hex.EncodeToString(publicKeyBytes))
	publicKeyStr := base64.StdEncoding.EncodeToString(publicKeyPem)
	fmt.Println("PUBLIC KEY LEN: ", len(publicKeyBytes)) // 550
	fmt.Println("PUBLIC KEY :", publicKeyStr)

	privateKeyBytes, err := x509.MarshalPKCS8PrivateKey(privateKey)
	if err != nil {
		panic(err)
	}
	privateKeyPem := pem.EncodeToMemory(&pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: privateKeyBytes,
	})
	//fmt.Println(hex.EncodeToString(publicKeyBytes))
	privateKeyStr := base64.StdEncoding.EncodeToString(privateKeyPem)
	fmt.Println("PRIVATE KEY LEN: ", len(privateKeyBytes)) // 2373-2375
	fmt.Println("PRIVATE KEY: ", privateKeyStr)

	//根据公钥加密
	encryptedBytes, err := rsa.EncryptOAEP(
		sha256.New(),
		rand.Reader,
		&publicKey,
		[]byte("测试哈哈哈"), //需要加密的字符串
		nil)
	if err != nil {
		panic(err)
	}
	fmt.Println("encrypted message: ", hex.EncodeToString(encryptedBytes))
	//根据私钥解密
	decryptedBytes, err := privateKey.Decrypt(nil, encryptedBytes, &rsa.OAEPOptions{Hash: crypto.SHA256})
	if err != nil {
		panic(err)
	}
	fmt.Println("decrypted message: ", string(decryptedBytes))

	privateKey2, err := rsa.GenerateKey(strings.NewReader("d03278e13c4559ced6dce74f0eab841c90fa6e2d8958d7ee3ad626b62eb7a0528433f0871423586da15e6312f7ccb510aa5a886e4533dbd99944ebcec5211431777ec1dcbc316a05c4312489c0d28cb5b54dcc9ba1d8314945508acacdd66769e011fde42b213c85112989d4e4c13fbe703c07b10740ac335c284772dc5a0784ca65476eab7b135d001a9b8bea5ad6409462abb0c0ac0085d5b8b8379b147269bda56f7adb306aeedc70322c6d573943e68067a4a6eb3895197d346ef931a112dec01fda1de5d8e5bae70a87297c1c3d2f89094c0a9700ad0d930aa8469680253385c1be4dda1023f68b66bc5179d4303a2cfab6e6f9d1b6c3a063fc58ce1182626a5e47da0460f9b4fea041f528d3b77f4f896b9c51c82d72ec6e315b4521332bd27538c1a024ad7089c2cc1567ba3f2b2271f21c0b3c559d68b51d0052c2d37b25c30a388fde3cb84579f835e70542843d5522e72dbd9ba6f72febbdc84aeb2b1276cf670b843b0714d6878be79ae76c7eb88a68f918ad6b63be58e0815a06872f04e5593478a71dda191be1cae1e063c2cca42d8e05cb0bb722266476a31cc1415ebfcf5bdd8aeb093430d2605b3922bd64f6a72f1e645175798e775ede049767392422e02a49601bfd7499dceafed9bc0f5946768e8183a3563dab37f93c2791f93db12139e261a71819d0586bb3ae6bd451c289250407b49334b3d644e7"), 128)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(privateKey2)
}
