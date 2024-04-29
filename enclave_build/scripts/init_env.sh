#!/bin/bash

# ami id:   ami-0dc60db7bddbc88b8
# ami name: Bansir AWS Nitro Enclaves Developer AMI ver-2021.07.22-733b3abf-c58e-4e14-9423-e82a94bd2950
# 推荐 EC2: c5.xlarge


yum install -y make zlib zlib-devel gcc-c++ libtool openssl openssl-devel vim htop net-tools pcre pcre-devel nginx mariadb105-server* redis* cronie ncurses-devel byacc flex java-1.8.0-amazon-corretto

amazon-linux-extras install golang1.19 -y

echo "install iftop"
mkdir iftop && cd iftop
wget http://www.tcpdump.org/release/libpcap-1.9.1.tar.gz
tar -xzvf libpcap-1.9.1.tar.gz && cd libpcap-1.9.1 && ./configure && make -j10 &&  make install && cd ../
wget http://www.ex-parrot.com/~pdw/iftop/download/iftop-0.17.tar.gz
tar -xzvf iftop-0.17.tar.gz && cd iftop-0.17 && ./configure && make -j10 && make install && cd ../
cd ..
rm -rf iftop

# 禁用 history
echo "export HISTSIZE=0" >> /etc/profile
# source /etc/profile
