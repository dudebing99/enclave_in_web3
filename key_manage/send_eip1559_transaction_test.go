package key_manage

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"github.com/dudebing99/go-ethereum/common"
	"github.com/dudebing99/go-ethereum/core/types"
	"github.com/dudebing99/go-ethereum/crypto"
	"github.com/dudebing99/go-ethereum/ethclient"
	"log"
	"math/big"
	"testing"
)

func TestTransfer(t *testing.T) {
	client, err := ethclient.Dial("http://54.151.136.132:8545")
	if err != nil {
		t.Fatal(err)
		log.Fatal(err)
	}

	privateKey, err := crypto.HexToECDSA("aead75071f4a9437df36d08acdcbb78b8dca55d02f0631f33f72274e9ee45a98")
	if err != nil {
		log.Fatal(err)
	}

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("error casting public key to ECDSA")
	}

	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)

	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		log.Fatal(err)
	}

	value := big.NewInt(200000000000000000) // in wei (0.2 DEMON)
	gasLimit := uint64(30000)               // in units
	tip := big.NewInt(2000000000)           // maxPriorityFeePerGas = 2 Gwei
	feeCap := big.NewInt(20000000000)       // maxFeePerGas = 20 Gwei
	if err != nil {
		log.Fatal(err)
	}

	toAddress := common.HexToAddress("0xCb75C706a45fefF971359F53dF7DD6dF47a41013")
	var data []byte = []byte("666")

	chainID, err := client.NetworkID(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	tx := types.NewTx(&types.DynamicFeeTx{
		ChainID:   chainID,
		Nonce:     nonce,
		GasFeeCap: feeCap,
		GasTipCap: tip,
		Gas:       gasLimit,
		To:        &toAddress,
		Value:     value,
		Data:      data,
	})

	signedTx, err := types.SignTx(tx, types.LatestSignerForChainID(chainID), privateKey)

	if err != nil {
		log.Fatal(err)
	}

	err = client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("transaction hash: ", signedTx.Hash().Hex())
}
