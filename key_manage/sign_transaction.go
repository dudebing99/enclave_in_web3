package key_manage

import (
	"bytes"
	"encoding/hex"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rlp"
	"math/big"
)

func DecodeRawTransaction(rawTx string) (tx *types.Transaction, err error) {
	tx = new(types.Transaction)
	rawTxBytes, err := hex.DecodeString(rawTx)
	if err != nil {
		return tx, err
	}

	err = tx.UnmarshalBinary(rawTxBytes)
	if err != nil {
		return tx, err
	}

	//txJson, _ := tx.MarshalJSON()
	//fmt.Println("tx:", string(txJson))
	return tx, err
}

func DecodeRLPRawTransaction(rawTx string) (tx *types.Transaction, err error) {
	tx = new(types.Transaction)
	rawTxBytes, err := hex.DecodeString(rawTx)
	if err != nil {
		return tx, err
	}

	rlpStream := rlp.NewStream(bytes.NewBuffer(rawTxBytes), 0)
	if err = tx.DecodeRLP(rlpStream); err != nil {
		return tx, err
	}

	//txJson, _ := tx.MarshalJSON()
	//fmt.Println("tx:", string(txJson))
	return tx, err
}

func SignTransaction(rawTx string, chainId *big.Int, keyBytes string) (signedRawTxStr string, err error) {
	tx, err := DecodeRLPRawTransaction(rawTx)
	if err != nil {
		return signedRawTxStr, err
	}

	privateKey, err := crypto.HexToECDSA(keyBytes)
	if err != nil {
		return signedRawTxStr, err
	}

	// NewCancunSigner returns a signer that accepts
	// - EIP-4844 blob transactions
	// - EIP-1559 dynamic fee transactions
	// - EIP-2930 access list transactions,
	// - EIP-155 replay protected transactions, and
	// - legacy Homestead transactions.
	signedRawTx, err := types.SignTx(tx, types.NewEIP155Signer(chainId), privateKey)
	if err != nil {
		return signedRawTxStr, err
	}

	signedRawTxBytes := new(bytes.Buffer)
	err = signedRawTx.EncodeRLP(signedRawTxBytes)
	if err != nil {
		return signedRawTxStr, err
	}

	return hex.EncodeToString(signedRawTxBytes.Bytes()), nil
}

func SignTransaction2(rawTx string, chainId *big.Int, keyBytes string) (signedRawTxStr string, err error) {
	tx, err := DecodeRawTransaction(rawTx)
	if err != nil {
		return signedRawTxStr, err
	}

	privateKey, err := crypto.HexToECDSA(keyBytes)
	if err != nil {
		return signedRawTxStr, err
	}

	// NewCancunSigner returns a signer that accepts
	// - EIP-4844 blob transactions
	// - EIP-1559 dynamic fee transactions
	// - EIP-2930 access list transactions,
	// - EIP-155 replay protected transactions, and
	// - legacy Homestead transactions.
	signedRawTx, err := types.SignTx(tx, types.NewEIP155Signer(chainId), privateKey)
	if err != nil {
		return signedRawTxStr, err
	}

	signedRawTxBytes, err := signedRawTx.MarshalBinary()
	if err != nil {
		return signedRawTxStr, err
	}

	return hex.EncodeToString(signedRawTxBytes), nil
}
