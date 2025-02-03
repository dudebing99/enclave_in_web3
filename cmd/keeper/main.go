package main

import (
	"bufio"
	"enclave_in_web3/config"
	"enclave_in_web3/dtos"
	"enclave_in_web3/key_manage"
	"enclave_in_web3/utils"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/mdlayher/vsock"
	"golang.org/x/sys/unix"
	"io"
	"log"
	"math/big"
	"net"
	"strings"
)

func process(conn net.Conn, keeper *key_manage.Keeper) {
	defer conn.Close()

	reader := bufio.NewReader(conn)
	msgType, msg, err := utils.Decode(reader)
	if err == io.EOF {
		log.Println("read eof")
		return
	}
	if err != nil {
		log.Println("read error: ", err)
		return
	}

	if msgType == utils.SetEncryptionSeedReq {
		// 不打印设置种子消息，防止泄露种子
		log.Println("received data from client, type: ", msgType, ", length: ", len(msg))
	} else {
		log.Println("received data from client: ", string(msg), ", type: ", msgType, ", length: ", len(msg))
	}

	var internalError dtos.InternalError
	switch msgType {
	case utils.GenerateKeyReq:
		var req dtos.GenerateKeyReq
		reqJson := msg
		_ = json.Unmarshal(reqJson, &req)
		showPrivateKey := req.ShowPrivateKey
		keyId := req.KeyId

		enclaveManagedKey := keeper.Generate()
		if len(keyId) != 0 {
			enclaveManagedKey.KeyId = keyId
		}
		// 添加私钥
		err := keeper.AddKey(enclaveManagedKey, false)
		if err != nil {
			log.Println("add key error: ", err)
			// 异常处理
			internalError.ErrorMsg = fmt.Sprint("add key error: ", err)
			goto InternalError
		}

		rsp := dtos.GenerateKeyRsp{
			KeyId:      enclaveManagedKey.KeyId,
			Address:    enclaveManagedKey.Address,
			PrivateKey: "",
		}
		if showPrivateKey {
			rsp.PrivateKey = enclaveManagedKey.PrivateKey
		}

		// 写数据
		rspJson, _ := json.Marshal(rsp)
		log.Println("try to send data to client, type: ", utils.GenerateKeyRsp, ", length: ", len(rspJson))
		encoded, _ := utils.Encode(utils.GenerateKeyRsp, rspJson)

		_, err = conn.Write(encoded)
		if err != nil {
			log.Println("write error: ", err)
		}

		return

	case utils.AddKeyReq:
		var req dtos.AddKeyReq
		reqJson := msg
		_ = json.Unmarshal(reqJson, &req)
		keyId := req.KeyId
		address := req.Address
		privateKey := req.PrivateKey

		err := key_manage.Validate(privateKey, address)
		if err != nil {
			log.Println("validate error: ", err)
			// 异常处理
			internalError.ErrorMsg = fmt.Sprint("validate error: ", err)
			goto InternalError
		}

		// 添加私钥
		enclaveManagedKey := key_manage.EnclaveManagedKey{
			KeyId:      keyId,
			Address:    address,
			PrivateKey: privateKey,
		}
		err = keeper.AddKey(enclaveManagedKey, false)
		if err != nil {
			log.Println("add key error: ", err)
			// 异常处理
			internalError.ErrorMsg = fmt.Sprint("add key error: ", err)
			goto InternalError
		}

		rsp := dtos.AddKeyRsp{
			//
		}

		// 写数据
		rspJson, _ := json.Marshal(rsp)
		log.Println("try to send data to client, type: ", utils.AddKeyRsp, ", length: ", len(rspJson))
		encoded, _ := utils.Encode(utils.AddKeyRsp, rspJson)

		_, err = conn.Write(encoded)
		if err != nil {
			log.Println("write error: ", err)
		}

		return

	case utils.SetEncryptionSeedReq:
		var req dtos.SetEncryptionSeedReq
		reqJson := msg
		_ = json.Unmarshal(reqJson, &req)
		seed := req.EncryptionSeed

		err := keeper.SetEncryptionSeed(seed)
		if err != nil {
			log.Println("set encryption seed error: ", err)
			// 异常处理
			internalError.ErrorMsg = fmt.Sprint("set encryption seed error: ", err)
			goto InternalError
		}
		rsp := dtos.SetEncryptionSeedRsp{
			//
		}
		// 写数据
		rspJson, _ := json.Marshal(rsp)
		log.Println("try to send data to client, type: ", utils.SetEncryptionSeedRsp, ", length: ", len(rspJson))
		encoded, _ := utils.Encode(utils.SetEncryptionSeedRsp, rspJson)

		_, err = conn.Write(encoded)
		if err != nil {
			log.Println("write error: ", err)
		}

		return

	case utils.GenerateAddressReq:
		var req dtos.GenerateAddressReq
		reqJson := msg
		_ = json.Unmarshal(reqJson, &req)

		number := req.Number
		rsp := make(dtos.GenerateAddressRsp, 0)

		for i := 0; i < int(number); i++ {
			enclaveManagedKey := keeper.Generate()
			// 添加私钥
			err := keeper.AddKey(enclaveManagedKey, false)
			if err != nil {
				log.Println("add key error: ", err)
				// 异常处理
				internalError.ErrorMsg = fmt.Sprint("add key error: ", err)
				goto InternalError
			}

			encryptedKey, err := keeper.EncryptPrivateKey(enclaveManagedKey.PrivateKey)
			if err != nil {
				log.Println("encrypt key error: ", err)
				// 异常处理
				internalError.ErrorMsg = fmt.Sprint("encrypt key error: ", err)
				goto InternalError
			}

			rsp = append(rsp, dtos.Address{
				KeyId:      enclaveManagedKey.KeyId,
				Address:    enclaveManagedKey.Address,
				PrivateKey: encryptedKey,
			})
		}

		// 写数据
		rspJson, _ := json.Marshal(rsp)
		log.Println("try to send data to client, type: ", utils.GenerateAddressRsp, ", length: ", len(rspJson))
		encoded, _ := utils.Encode(utils.GenerateAddressRsp, rspJson)

		_, err = conn.Write(encoded)
		if err != nil {
			log.Println("write error: ", err)
		}

		return

	case utils.SignMessageReq:
		var req dtos.SignMessageReq
		reqJson := msg
		_ = json.Unmarshal(reqJson, &req)
		keyId := req.KeyId
		message := req.Message
		needToHash := req.NeedToHash

		privateKey, err := keeper.Get(keyId)
		if err != nil {
			log.Println("get key error: ", err)
			// 异常处理
			internalError.ErrorMsg = err.Error()
			goto InternalError
		}

		// 去掉前缀 0x
		message = strings.TrimPrefix(message, "0x")
		decodedMessage, err := hex.DecodeString(message)
		if err != nil {
			log.Println("decode message error: ", err)
			// 异常处理
			internalError.ErrorMsg = fmt.Sprint("decode message error: ", err)
			goto InternalError
		}
		signature, err := key_manage.Sign(privateKey, decodedMessage, needToHash)
		if err != nil {
			log.Println("sign message error: ", err)
			// 异常处理
			internalError.ErrorMsg = fmt.Sprint("sign message error: ", err)
			goto InternalError
		}
		rsp := dtos.SignMessageRsp{
			Signature: signature,
		}
		// 写数据
		rspJson, _ := json.Marshal(rsp)
		log.Println("try to send data to client, type: ", utils.SignMessageRsp, ", length: ", len(rspJson))
		encoded, _ := utils.Encode(utils.SignMessageRsp, rspJson)

		_, err = conn.Write(encoded)
		if err != nil {
			log.Println("write error: ", err)
		}

		return

	case utils.AddAddressReq:
		var req dtos.AddAddressReq
		reqJson := msg
		_ = json.Unmarshal(reqJson, &req)
		keyId := req.KeyId
		address := req.Address
		encodedEncryptedPrivateKey := req.PrivateKey // 加密后的私钥

		// 需要先解密
		err = keeper.IsReady()
		if err != nil {
			log.Println("not ready: ", err)
			// 异常处理
			internalError.ErrorMsg = fmt.Sprint("not ready: ", err)
			goto InternalError
		}

		privateKey, err := keeper.DecryptPrivateKey(encodedEncryptedPrivateKey)
		if err != nil {
			log.Println("decrypt private key error: ", err)
			// 异常处理
			internalError.ErrorMsg = fmt.Sprint("decrypt private key error: ", err)
			goto InternalError
		}

		err = key_manage.Validate(privateKey, address)
		if err != nil {
			log.Println("validate error: ", err)
			// 异常处理
			internalError.ErrorMsg = fmt.Sprint("validate error: ", err)
			goto InternalError
		}

		// 添加私钥
		enclaveManagedKey := key_manage.EnclaveManagedKey{
			KeyId:      keyId,
			Address:    address,
			PrivateKey: privateKey,
		}
		err = keeper.AddKey(enclaveManagedKey, true)
		if err != nil {
			log.Println("add key error: ", err)
			// 异常处理
			internalError.ErrorMsg = fmt.Sprint("add key error: ", err)
			goto InternalError
		}

		rsp := dtos.AddAddressRsp{
			//
		}

		// 写数据
		rspJson, _ := json.Marshal(rsp)
		log.Println("try to send data to client, type: ", utils.AddAddressRsp, ", length: ", len(rspJson))
		encoded, _ := utils.Encode(utils.AddAddressRsp, rspJson)

		_, err = conn.Write(encoded)
		if err != nil {
			log.Println("write error: ", err)
		}

		return

	case utils.SignTransactionReq:
		var req dtos.SignTransactionReq
		reqJson := msg
		_ = json.Unmarshal(reqJson, &req)
		keyId := req.KeyId
		chainId := req.ChainId
		rawTx := req.RawTx

		privateKey, err := keeper.Get(keyId)
		if err != nil {
			log.Println("get key error: ", err)
			// 异常处理
			internalError.ErrorMsg = err.Error()
			goto InternalError
		}

		// 去掉前缀 0x
		rawTx = strings.TrimPrefix(rawTx, "0x")
		signedRawTx, err := key_manage.SignTransaction(rawTx, big.NewInt(int64(chainId)), privateKey)
		if err != nil {
			log.Println("sign tx error: ", err)
			// 异常处理
			internalError.ErrorMsg = fmt.Sprint("sign tx error: ", err)
			goto InternalError
		}
		rsp := dtos.SignTransactionRsp{
			SignedRawTx: signedRawTx,
		}
		// 写数据
		rspJson, _ := json.Marshal(rsp)
		log.Println("try to send data to client, type: ", utils.SignTransactionRsp, ", length: ", len(rspJson))
		encoded, _ := utils.Encode(utils.SignTransactionRsp, rspJson)

		_, err = conn.Write(encoded)
		if err != nil {
			log.Println("write error: ", err)
		}

		return

	default:
		log.Println("unknown message type!")
		internalError.ErrorMsg = fmt.Sprint("unknown message type")
		goto InternalError
	}

InternalError:
	// 写数据
	rspJson, _ := json.Marshal(internalError)
	log.Println("try to send internal error to client, type: ", utils.InternalErrorType, ", length: ", len(rspJson))
	encoded, _ := utils.Encode(utils.InternalErrorType, rspJson)

	_, err = conn.Write(encoded)
	if err != nil {
		log.Println("write error: ", err)
	}
}

func main() {
	fmt.Println("Starting enclave keeper ...")

	// 初始化配置文件
	config.InitConfig()

	keeper := key_manage.NewKeeper()

	// 测试网测试环境预置签名私钥
	err := keeper.AddKey(
		key_manage.EnclaveManagedKey{
			KeyId:      "f47ac10b-58cc-0372-8567-0e02b2c3d479",
			Address:    "0xCb75C706a45fefF971359F53dF7DD6dF47a41013",
			PrivateKey: "aead75071f4a9437df36d08acdcbb78b8dca55d02f0631f33f72274e9ee45a98",
		}, false)
	utils.CheckError(err)

	// 主网测试环境预置签名私钥
	err = keeper.AddKey(key_manage.EnclaveManagedKey{
		KeyId:      "6ddcd9c1-7a6a-42b0-8641-4311b4cb98b4",
		Address:    "0xE7c441409A79E8Eec7489de81697b3fE44281182",
		PrivateKey: "dfd5b91a521e985eef6d2b46cd0b170f72b0315c741b1e7389e1e4493c9e4f6f",
	}, false)
	utils.CheckError(err)

	// Listen for VM sockets connections on port 1024.
	l, err := vsock.ListenContextID(unix.VMADDR_CID_ANY, 1024, nil)
	utils.CheckError(err)

	defer func(l *vsock.Listener) {
		err := l.Close()
		utils.CheckError(err)
	}(l)

	// Accept a single connection.
	c, err := l.Accept()
	utils.CheckError(err)

	defer func(c net.Conn) {
		err := c.Close()
		utils.CheckError(err)
	}(c)

	for {
		conn, err := l.Accept()
		utils.CheckError(err)

		go process(conn, keeper)
	}
}
