package utils

import (
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/google/uuid"
)

func IsValidKeyId(keyId string) bool {
	_, err := uuid.Parse(keyId)
	return err == nil
}

func IsValidEthereumAddress(address string) bool {
	return common.IsHexAddress(address)
}

func IsValidPrivateKey(privateKey string) bool {
	privateKeyBytes := make([]byte, 32)
	_, err := fmt.Sscanf(privateKey, "%64x", &privateKeyBytes)

	if err != nil {
		return false
	}

	_, err = crypto.ToECDSA(privateKeyBytes)
	return err == nil
}
