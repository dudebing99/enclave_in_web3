package dtos

type SignMessageReq struct {
	KeyId      string `json:"key_id"`
	Message    string `json:"message"`
	NeedToHash bool   `json:"need_to_hash"`
}

type SignMessageRsp struct {
	Signature string `json:"signature"`
}

type SignTransactionReq struct {
	KeyId   string `json:"key_id"`
	ChainId uint32 `json:"chain_id"`
	RawTx   string `json:"raw_tx"`
}

type SignTransactionRsp struct {
	SignedRawTx string `json:"signed_raw_tx"`
}
