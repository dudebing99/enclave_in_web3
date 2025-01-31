package key_manage

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"golang.org/x/crypto/sha3"
	"log"
	"math/big"
	"testing"
)

func TestTransferErc20(t *testing.T) {
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
		log.Fatal("cannot assert type: publicKey is not of type *ecdsa.PublicKey")
	}

	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("fromAddress: ", fromAddress, ", nonce: ", nonce)

	value := big.NewInt(0) // in wei (0 eth)
	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("gasPrice: ", gasPrice)

	toAddress := common.HexToAddress("0xCb75C706a45fefF971359F53dF7DD6dF47a41013")
	tokenAddress := common.HexToAddress("0xCdA67E179f6c6773f2529e2130bDbf472b3F16Ad") // GM

	transferFnSignature := []byte("transfer(address,uint256)")
	hash := sha3.NewLegacyKeccak256()
	hash.Write(transferFnSignature)
	methodID := hash.Sum(nil)[:4]
	fmt.Println(hexutil.Encode(methodID)) // 0xa9059cbb

	paddedAddress := common.LeftPadBytes(toAddress.Bytes(), 32)
	fmt.Println(hexutil.Encode(paddedAddress))

	amount := new(big.Int)
	amount.SetString("2", 10)

	paddedAmount := common.LeftPadBytes(amount.Bytes(), 32)
	fmt.Println(hexutil.Encode(paddedAmount))

	var data []byte
	data = append(data, methodID...)
	data = append(data, paddedAddress...)
	data = append(data, paddedAmount...)

	//gasLimit, err := client.EstimateGas(context.Background(), ethereum.CallMsg{
	//	To:   &tokenAddress,
	//	Data: data,
	//})
	//if err != nil {
	//	log.Fatal("", err)
	//}
	//fmt.Println("gasLimit: ", gasLimit)

	tx := types.NewTransaction(nonce, tokenAddress, value, 66666, gasPrice, data)

	chainID, err := client.NetworkID(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("chainId: ", chainID)

	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
	if err != nil {
		log.Fatal(err)
	}

	err = client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("transaction hash: ", signedTx.Hash().Hex())
}
