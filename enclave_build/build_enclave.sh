#/bin/bash

#nitro-cli terminate-enclave --enclave-id $(nitro-cli describe-enclaves | jq -r ".[0].EnclaveID")

cp ../bin/enclave-keeper* enclave_keeper
cp ../conf/application.yml application.yml

# 删除无标签且无引用的镜像
docker image prune

docker build ./ -t enclave_keeper
nitro-cli build-enclave --docker-uri enclave_keeper:latest --output-file enclave_keeper.eif
#nitro-cli run-enclave --cpu-count 2 --memory 3000 --enclave-cid 22 --eif-path enclave_keeper.eif --debug-mode
nitro-cli run-enclave --cpu-count 2 --memory 3000 --enclave-cid 22 --eif-path enclave_keeper.eif
#nitro-cli describe-enclaves
#nitro-cli console --enclave-id $(nitro-cli describe-enclaves | jq -r ".[0].EnclaveID")
