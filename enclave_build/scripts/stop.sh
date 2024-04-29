#!/bin/bash

EnclaveID=`nitro-cli describe-enclaves|grep EnclaveID|awk -F '"' '{print $4}'`

echo "EnclaveID: "$EnclaveID

nitro-cli terminate-enclave --enclave-id $EnclaveID
