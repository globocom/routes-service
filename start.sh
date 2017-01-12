#!/bin/bash

# Exports PATH 
export PATH=$PATH:/usr/local/go/bin
# Exports GOPATH
export GOPATH=$HOME
#Start gobgpd and webserver
$HOME/bin/gobgpd & go run main.go & 
