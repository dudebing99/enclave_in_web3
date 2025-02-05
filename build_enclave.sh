#!/bin/bash

#nitro-cli terminate-enclave --enclave-id $(nitro-cli describe-enclaves | jq -r ".[0].EnclaveID")

docker build ./ -t enclave_keeper:latest

docker image prune

nitro-cli build-enclave --docker-uri enclave_keeper:latest --output-file enclave_keeper.eif
#nitro-cli run-enclave --cpu-count 2 --memory 3000 --enclave-cid 22 --eif-path enclave_keeper.eif --debug-mode
nitro-cli run-enclave --cpu-count 2 --memory 3000 --enclave-cid 22 --eif-path enclave_keeper.eif
#nitro-cli describe-enclaves
#nitro-cli console --enclave-id $(nitro-cli describe-enclaves | jq -r ".[0].EnclaveID")
