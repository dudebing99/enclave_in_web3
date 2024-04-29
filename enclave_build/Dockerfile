# Copyright 2020 Amazon.com, Inc. or its affiliates. All Rights Reserved.
# SPDX-License-Identifier: Apache-2.0

FROM centos:7
MAINTAINER kevin <admin@gmail.com>

ENV HELLO="Hello from the enclave side!"
COPY enclave_keeper /bin/enclave_keeper
RUN mkdir /etc/conf
COPY application.yml /etc/conf

# ./enclave_keeper -config /etc/conf/application.yml
CMD ["/bin/enclave_keeper", "-config", "/etc/conf/application.yml"]
