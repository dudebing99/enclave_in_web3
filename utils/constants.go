package utils

const (
	KeyIdLength      = 36
	AddressLength    = 42
	PrivateKeyLength = 64
)

// DefaultMaxPrivateKeys
// 需要考虑进程运行所需内存，实际场景下，如果大型交易所，管理海量地址，可以设置多个 Enclave 实例，业务层
// 在保存用户（地址）数据表中增加关联的 Enclave 实例的 Id，业务调用时根据实例 Id 路由到不同的 Enclave 实例
// 即可。 例如，管理 100 万个密钥所需内存不超过 300 MB
const DefaultMaxPrivateKeys uint32 = 1000000
