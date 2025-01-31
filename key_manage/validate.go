package key_manage

import (
	"errors"
	"github.com/ethereum/go-ethereum/crypto"
)

func Validate(keyBytes string, address string) error {
	privateKey, err := crypto.HexToECDSA(keyBytes)
	if err != nil {
		return err
	}

	addr := crypto.PubkeyToAddress(privateKey.PublicKey).Hex()
	if addr != address {
		return errors.New("address and private key mismatched")
	}
	return nil
}
