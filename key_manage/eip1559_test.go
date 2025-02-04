package key_manage

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"math/big"
	"testing"
)

func TestGenerateEip1559Transaction(t *testing.T) {
	value := big.NewInt(1)  // in wei
	gasLimit := uint64(1)   // in units
	tip := big.NewInt(1)    // maxPriorityFeePerGas = 2 Gwei
	feeCap := big.NewInt(1) // maxFeePerGas = 20 Gwei
	chainId := big.NewInt(1)
	toAddress := common.HexToAddress("0xCb75C706a45fefF971359F53dF7DD6dF47a41013")
	nonce := uint64(1)
	//data := []byte("crazy eip1559")
	tx := types.NewTx(&types.DynamicFeeTx{
		ChainID:   chainId,
		Nonce:     nonce,
		GasFeeCap: feeCap,
		GasTipCap: tip,
		Gas:       gasLimit,
		To:        &toAddress,
		Value:     value,
		Data:      nil,
		//AccessList: nil,
	})

	fmt.Println("type: ", uint32(tx.Type()))
	rawTxBytes := new(bytes.Buffer)
	err := tx.EncodeRLP(rawTxBytes)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("rlp raw tx: ", hex.EncodeToString(rawTxBytes.Bytes()))

	rawTxBytes2, err := tx.MarshalBinary()
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("raw tx: ", hex.EncodeToString(rawTxBytes2))
}
