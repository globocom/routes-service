#!/bin/bash

apt-get update
apt-get install git -y

wget https://storage.googleapis.com/golang/go1.7.4.linux-amd64.tar.gz

tar -C /usr/local -xzf go$VERSION.$OS-$ARCH.tar.gz

export PATH=$PATH:/usr/local/go/bin

mkdir $HOME/work
export GOPATH=$HOME/work

cd work

go get github.com/osrg/gobgp/gobgpd
go get github.com/osrg/gobgp/gobgp