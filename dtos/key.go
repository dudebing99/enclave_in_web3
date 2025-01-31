package dtos

type GenerateKeyReq struct {
	ShowPrivateKey bool   `json:"show_private_key"`
	KeyId          string `json:"key_id"` // 可以指定 key id，便于业务层不用修改 key id 即可使用新私钥
}

type GenerateKeyRsp struct {
	KeyId      string `json:"key_id"`
	Address    string `json:"address"`
	PrivateKey string `json:"private_key"`
}

type AddKeyReq struct {
	KeyId      string `json:"key_id"`
	Address    string `json:"address"`
	PrivateKey string `json:"private_key"`
}

type AddKeyRsp struct {
	//
}

type SetEncryptionSeedReq struct {
	EncryptionSeed string `json:"encryption_seed"`
}

type SetEncryptionSeedRsp struct {
	//
}

type GenerateAddressReq struct {
	ShowPrivateKey bool   `json:"show_private_key"` // 加密后的私钥
	Number         uint32 `json:"number"`           // 支持批量生成地址
}

type Address struct {
	KeyId      string `json:"key_id"`
	Address    string `json:"address"`
	PrivateKey string `json:"private_key"` // 加密后的私钥
}

type GenerateAddressRsp []Address

type AddAddressReq struct {
	KeyId      string `json:"key_id"`
	Address    string `json:"address"`
	PrivateKey string `json:"private_key"` // 加密后的私钥
}

type AddAddressRsp struct {
	//
}
