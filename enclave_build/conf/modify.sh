#/bin/bash

cp -f ./allocator.yaml /etc/nitro_enclaves/allocator.yaml

systemctl restart nitro-enclaves-allocator.service
