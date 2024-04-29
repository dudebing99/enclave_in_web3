const Web3 = require('web3');
const Tx = require('ethereumjs-tx');
const Http = require('http');

const web3 = new Web3('https://endpoints.omniatech.io/v1/bsc/testnet/public');
//const web3 = new Web3('https://flashy-solitary-diagram.bsc.quiknode.pro/429d3370c64df671c118a5fd847096d012843e88/')

function signTransaction(data, callback) {
    return new Promise(resolve => {
        const opt = {
            hostname: '13.228.243.51',
            port: 10000,
            path: '/8ed8d7fe15437d09aeba8b757cc14cdc/api/sign/transaction',
            method: 'POST',
            data: data,
            headers: {
                'Content-Type': 'application/json; charset=UTF-8',
                'Content-Length': data.length
            }
        };

        let body = '';
        let req = Http.request(opt, function (res) {
            // res.statusCode == 200;
            res.setEncoding('utf8');
            res.on('data', function (chunk) {
                body += chunk;
            }).on('end', function () {
                resolve(callback(body));
            });
        }).on('error', function (err) {
            console.log('error: ', err.message);
        });

        req.write(data);
        req.end();
    })
}

// Create an async function so I can use the 'await' keyword to wait for things to finish
async function main() {
    console.log(`web3 version: ${web3.version}`)

    let addr = '0xCb75C706a45fefF971359F53dF7DD6dF47a41013';
    // let bal = await web3.eth.getBalance(addr);
    // let balanceInEther = web3.utils.fromWei(bal, 'ether')
    let nonce = await web3.eth.getTransactionCount(addr);
    console.log(`current nonce: ${nonce}`);

    let rawTransaction = {
        'from': addr,
        'nonce': '0x' + nonce.toString(16),
        'gasPrice': '0x12a05f200', // 5 GWei
        'gasLimit': '0x250CA',
        'to': addr,
        'value': '0x666',
        'data': '666',
        //'chainId': 0x38 // mainnet 56
        'chainId': 0x61 // testnet 97
    };

    let tx = new Tx(rawTransaction);
    let serializedTx = tx.serialize();
    let rawTx = serializedTx.toString('hex')
    console.log(`unsigned tx: ${rawTx}`);


    let req = {
        "key_id": "f47ac10b-58cc-0372-8567-0e02b2c3d479",
        "chain_id": 97,
        "raw_tx": rawTx
    }

    let signedRawTx = await signTransaction(JSON.stringify(req), function (res) {
        // console.log(res);
        let obj = JSON.parse(res);
        let signedRawTx = obj.data.signed_raw_tx;
        return signedRawTx;
    });

    console.log(`signed tx: ${signedRawTx}`);
    let receipt = await web3.eth.sendSignedTransaction("0x" + signedRawTx);
    console.log(`${JSON.stringify(receipt, null, '\t')}`);
}

main();
