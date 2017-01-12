#!/bin/bash

apt-get update
apt-get install git -y

# Download GO
cd /tmp/
wget https://storage.googleapis.com/golang/go1.7.4.linux-amd64.tar.gz

# Extracts Files
tar -C /usr/local -xzf go1.7.4.linux-amd64.tar.gz

# Exports PATH 
export PATH=$PATH:/usr/local/go/bin

# Exports GOPATH
export GOPATH=$HOME

# Install gobgp
cd /vagrant
go get github.com/osrg/gobgp/gobgpd
go get github.com/osrg/gobgp/gobgp

#Install echo
go get github.com/labstack/echo

#Start gobgpd and webserver
$HOME/bin/gobgpd & go run main.go & 
