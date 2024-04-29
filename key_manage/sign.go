package key_manage

import (
	"github.com/dudebing99/go-ethereum/accounts"
	"github.com/dudebing99/go-ethereum/common/hexutil"
	"github.com/dudebing99/go-ethereum/crypto"
)

// 签名中不同语言库存在差异, 解决方案参考资料如下:
// 1. golang 签名 v 始终为 0 或 1:	https://pkg.go.dev/github.com/ethereum/go-ethereum/crypto#Sign
// 2. 与 Java 不兼容, 解决方案:		https://ethereum.stackexchange.com/questions/45580/validating-go-ethereum-key-signature-with-ecrecover
// 3. 常用签名验证方案:				https://github.com/storyicon/sigverify

func Sign(keyBytes string, message []byte, needToHash bool) (string, error) {
	privateKey, err := crypto.HexToECDSA(keyBytes)
	var messageHash []byte
	if needToHash {
		messageHash = accounts.TextHash(message)
	} else {
		messageHash = message
	}

	signature, err := crypto.Sign(messageHash, privateKey)
	if err != nil {
		return "", err
	}

	// 始终为真
	if len(signature) > 64 && (uint8(signature[64]) == 0 || uint8(signature[64]) == 1) {
		signature[64] = signature[64] + 27
	}

	return hexutil.Encode(signature), nil
}
