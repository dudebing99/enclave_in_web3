package key_manage

import (
	"encoding/hex"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/cmd/utils"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/google/uuid"
)

type EnclaveManagedKey struct {
	KeyId      string
	Address    string
	PrivateKey string
}

func Generate() EnclaveManagedKey {
	privateKey, err := crypto.GenerateKey()
	if err != nil {
		utils.Fatalf("Failed to generate random private key: %v", err)
	}

	// Create the keyfile object with a random UUID.
	UUID, err := uuid.NewRandom()
	if err != nil {
		utils.Fatalf("Failed to generate random uuid: %v", err)
	}

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
