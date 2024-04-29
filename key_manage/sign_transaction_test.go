package key_manage

import (
	"fmt"
	"math/big"
	"testing"
)

func TestDecodeRLPTransaction(t *testing.T) {
	rawTx := "e8823f3085012a05f20082753094cb75c706a45feff971359f53df7dd6df47a41013823124830cfaef"

	tx, err := DecodeRLPRawTransaction(rawTx)
	if err != nil {
		t.Fatalf("%v", err)
	} else {
		txJson, _ := tx.MarshalJSON()
		fmt.Println(string(txJson))
	}

	keyBytes := "aead75071f4a9437df36d08acdcbb78b8dca55d02f0631f33f72274e9ee45a98"
	signedRawTxStr, err := SignTransaction(rawTx, big.NewInt(97), keyBytes)
	fmt.Println(signedRawTxStr)

	tx, err = DecodeRLPRawTransaction(signedRawTxStr)
	if err != nil {
		t.Fatalf("%v", err)
	} else {
		txJson, _ := tx.MarshalJSON()
		fmt.Println(string(txJson))
	}
}

func TestDecodeTransaction(t *testing.T) {
	rawTx := "02ee820617823f2c8504a817c80084b2d05e0082520894cb75c706a45feff971359f53df7dd6df47a4101382312480c0"
	//02 e0 0101010101 94 cb75c706a45feff971359f53df7dd6df47a41013 0180c0 808080 	-- Go
	//02 dd 0101010101 94 cb75c706a45feff971359f53df7dd6df47a41013 0180c0			-- Java
	tx, err := DecodeRawTransaction(rawTx)
	if err != nil {
		t.Fatalf("%v", err)
	} else {
		txJson, err := tx.MarshalJSON()
		if err != nil {
			t.Fatalf("%v", err)
		}
		fmt.Println(string(txJson))
	}

	keyBytes := "aead75071f4a9437df36d08acdcbb78b8dca55d02f0631f33f72274e9ee45a98"
	signedRawTxStr, err := SignTransaction2(rawTx, big.NewInt(97), keyBytes)
	if err != nil {
		t.Fatalf("%v", err)
	}

	tx, err = DecodeRawTransaction(signedRawTxStr)
	if err != nil {
		t.Fatalf("%v", err)
	} else {
		txJson, _ := tx.MarshalJSON()
		fmt.Println(string(txJson))
	}
}
