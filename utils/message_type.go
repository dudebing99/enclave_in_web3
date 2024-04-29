package utils

type MessageType uint32

const (
	UnknownMessageType MessageType = iota
	GenerateKeyReq
	GenerateKeyRsp
	SignMessageReq
	SignMessageRsp
	AddKeyReq
	AddKeyRsp
	SignTransactionReq
	SignTransactionRsp
	SetEncryptionSeedReq
	SetEncryptionSeedRsp
	GenerateAddressReq
	GenerateAddressRsp
	AddAddressReq
	AddAddressRsp

	InternalErrorType MessageType = 10000
)
