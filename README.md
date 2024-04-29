# enclave_client

## 接口

- 服务健康检测接口 `/api/health` `GET` `POST`

示例请求

```bash
curl http://localhost:10000/api/health -s|jq
{
  "data": null,
  "error_code": 0,
  "error_msg": "ok",
  "req_id": "660f89fa334fd4d"
}
```

- 生成私钥接口 `/api/key/generate` `POST`

示例请求

```bash
curl http://localhost:10000/api/key/generate -X POST -d '{}' -s|jq
{
  "data": {
    "key_id": "540c039e-b7ef-493f-bf66-7d6cdbcc8b34",
    "address": "0x76693f38c23dF786756ff9C55B0AcFFE245b680f",
    "private_key": ""
  },
  "error_code": 0,
  "error_msg": "ok",
  "req_id": "a4714da01c2705d"
}
```

> 测试模式下，可以返回生成的私钥，生产模式下不建议使用

```bash
curl http://localhost:10000/api/key/generate -X POST -d '{"show_private_key":true}' -s|jq
{
  "data": {
    "key_id": "540c039e-b7ef-493f-bf66-7d6cdbcc8b34",
    "address": "0x162E4cc18A66731571A244840645fd58048B308f",
    "private_key": "1db13f9c3a3c91c9949f26040d33504ccb08d71acff4f76782b897a9383b7dc4"
  },
  "error_code": 0,
  "error_msg": "ok",
  "req_id": "ddc161a32bf1958"
}
```

- 添加私钥接口 `/api/key/add` `POST`

示例请求

```bash
curl http://localhost:10000/api/key/add -X POST -d '{"key_id":"3bf2b6cb-debe-4c75-bb6e-38a35e818680","address":"0x5d1D0e4A7775BD50565af85b679E726648753bC5","private_key":"ac904c67de52249561ea562f3190a7c875346c697ad80d8b037fecf688959293"}' -s|jq
{
    "data": {},
    "error_code": 0,
    "error_msg": "ok",
    "req_id": "664f14586ce69ba"
}
```

- 签名消息接口 `/api/sign/message` `POST`

示例请求

```bash
curl http://localhost:10000/api/sign/message -X POST -d '{"key_id": "540c039e-b7ef-493f-bf66-7d6cdbcc8b34", "message": "a90abb62e8c7614e8f8af810083db2b4d1ddb47c2839457ea539fd232feaeed3"}' -s|jq
{
  "data": {
    "signature": "168d885cf992d7cc1b7831a65ae184b23f21548897477cf9b452243baa21bbf1216734d953a549c7b5014e33cfd6d84d4c0df1cae80c0725a12e884ac442a52e1c"
  },
  "error_code": 0,
  "error_msg": "ok",
  "req_id": "920d94cf279f8a9"
}
```

> 如果请求的签名私钥不存在，错误如下

```bash
curl http://localhost:10000/api/sign/message -X POST -d '{"key_id": "b16a2659-a332-4293-9604-a8d1b3f879a3", "message": "a90abb62e8c7614e8f8af810083db2b4d1ddb47c2839457ea539fd232feaeed3"}' -s|jq
{
  "data": null,
  "error_code": 2,
  "error_msg": "fetch key error: not found",
  "req_id": "97668ad762a6675"
}
```

- 签名交易接口 `/api/sign/transaction` `POST`

> 交易包括切不限于主币转账、代币转账、其他合约交易

示例请求

1.使用 Java 生成交易，以 BSC 测试网为例

```java
EthGetTransactionCount ethGetTransactionCount = web3j.ethGetTransactionCount(credentials.getAddress(), DefaultBlockParameterName.PENDING).send();
BigInteger nonce =  ethGetTransactionCount.getTransactionCount();
// 纯转账交易，data 可缺省
// RawTransaction rawTransaction = RawTransaction.createEtherTransaction(nonce, gasPrice, gasLimit, to, value);

// 如果转账交易携带数据或者本身为合约交易，data 自行先组装好
RawTransaction raw = RawTransaction.createTransaction(
        new BigInteger("16176"),
        new BigInteger("5000000000"),
        new BigInteger("30000"),
        "0xCb75C706a45fefF971359F53dF7DD6dF47a41013",
        new BigInteger("12580"),
        "crasy");
System.out.println(Numeric.toHexString(TransactionEncoder.encode(raw)));
//0xe8823f3085012a05f20082753094cb75c706a45feff971359f53df7dd6df47a41013823124830cfaef
```

2.离线签名

> BSC 主网：chain_id=65；测试网：chain_id=97

```bash
curl http://localhost:10000/api/sign/transaction -X POST -d '{"key_id": "f47ac10b-58cc-0372-8567-0e02b2c3d479", "chain_id": 97, "raw_tx": "e8823f3085012a05f20082753094cb75c706a45feff971359f53df7dd6df47a41013823124830cfaef"}' -s|jq
{
  "data": {
    "signed_raw_tx": "f86c823f3085012a05f20082753094cb75c706a45feff971359f53df7dd6df47a41013823124830cfaef81e5a090690d83f346b06e9301450f1923f7e89c425b166b45b26027be9465dfc0f540a032b242962b7028326cf906d1a88fe783466bff39035ef60f37260e369e26df15"
  },
  "error_code": 0,
  "error_msg": "ok",
  "req_id": "1af51811b72a5bb"
}
```

3.广播交易

> Java 业务层广播，将上一步返回的数据组装好："0x" + <signed_raw_tx>

```java
EthSendTransaction ethSendTransaction = web3j.ethSendRawTransaction(hexValue).sendAsync().get();
String transactionHash = ethSendTransaction.getTransactionHash();
// poll for transaction response via org.web3j.protocol.Web3j.ethGetTransactionReceipt(<txHash>)
```

> 此处以直接调用 RPC 广播为例

```bash
curl --location 'https://endpoints.omniatech.io/v1/bsc/testnet/public' \
--header 'Content-Type: application/json' \
--data '{
    "jsonrpc": "2.0",
    "method": "eth_sendRawTransaction",
    "params": [
        "0xf86c823f3085012a05f20082753094cb75c706a45feff971359f53df7dd6df47a41013823124830cfaef81e5a090690d83f346b06e9301450f1923f7e89c425b166b45b26027be9465dfc0f540a032b242962b7028326cf906d1a88fe783466bff39035ef60f37260e369e26df15"
    ],
    "id": 1
}'
```

对应，响应如下

```bash
{
    "jsonrpc": "2.0",
    "id": 1,
    "result": "0x015a09d95d1b06f56b8c157b0b3101b498f2fc9c4f8815b5d1b99716a4707ec8"
}
```

- 设置加密种子接口 `/api/key/set_encryption_seed` `POST`

示例请求

```bash
curl --location 'http://localhost:10000/8ed8d7fe15437d09aeba8b757cc14cdc/api/key/set_encryption_seed' \
--header 'Content-Type: application/json' \
--data '{"encryption_seed":"oXuVXyjjGR5YdmNxf&Qjut2b-6s*wnU!DOqPSAUW#(5k!rrYnn)uAc_DOL-VkENJU-"}'
```

- 生成地址接口 `/api/key/generate_address` `POST`

示例请求

```bash
curl --location 'http://localhost:10000/api/key/generate_address' \
--header 'Content-Type: application/json' \
--data '{}'
```

```bash
{
    "data": {
        "key_id": "6859c95b-92df-4058-9e7c-75192b369cb5",
        "address": "0x7Dd4680922A2600046f231e88F702515AE716b06",
        "private_key": "63cd11211f1aa0d29a14258acc1ce39d01ee22eecae61424613b87b555b6c9b152b653435be00c3109aa869a3e5fbd8ce0e93ef78a9ab27f0583e4863eadd5d736e4bc7f5b33f5b13eeea611e98677ac"
    },
    "error_code": 0,
    "error_msg": "ok",
    "req_id": "3c407421941b9c2"
}
```

# enclave_keeper

## 参考资料

- [vsock in Go] https://github.com/mdlayher/vsock
- [Installing the Nitro Enclaves CLI on Linux] https://docs.aws.amazon.com/enclaves/latest/user/nitro-enclave-cli-install.html
- [Hello enclave] https://docs.aws.amazon.com/enclaves/latest/user/getting-started.html
- [aws-nitro-enclaves-samples] https://github.com/aws/aws-nitro-enclaves-samples
- [部署一个 Nitro Enclave 示例环境，结合 KMS 实现私钥安全] https://github.com/hxhwing/Nitro-Enclave-Demo
- [基于 Nitro Enclave 构建安全的可信执行环境] https://aws.amazon.com/cn/blogs/china/build-a-secure-and-trusted-execution-environment-based-on-nitro-enclave/
- [kms tool] https://github.com/aws/aws-nitro-enclaves-sdk-c
