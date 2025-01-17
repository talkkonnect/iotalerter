#!/bin/bash

apt-get update
apt-get -y dist upgrade

## Add talkkonnect user to the system
adduser --disabled-password --disabled-login --gecos "" talkkonnect
usermod -a -G cdrom,audio,video,plugdev,users,dialout,dip,input,gpio talkkonnect

## Install the dependencies required for iotalerter
apt-get -y install git screen pkg-config

## Create the necessary directory structure under /home/talkkonnect/
cd /home/talkkonnect/
mkdir -p /home/talkkonnect/gocode
mkdir -p /home/talkkonnect/bin

## Create the log file
touch /var/log/iotalerter.log

# Check Latest of GOLANG 64 Bit Version for Raspberry Pi
GOLANG_LATEST_STABLE_VERSION=$(curl -s https://go.dev/VERSION?m=text | grep go)
cputype=`lscpu | grep Architecture | cut -d ":" -f 2 | sed 's/ //g'`
bitsize=`getconf LONG_BIT`

cd /usr/local

if [ $bitsize == '32' ]
then
echo "32 bit processor"
wget -nc https://go.dev/dl/$GOLANG_LATEST_STABLE_VERSION.linux-armv6l.tar.gz $GOLANG_LATEST_STABLE_VERSION.linux-armv6l.tar.gz
tar -zxvf /usr/local/$GOLANG_LATEST_STABLE_VERSION.linux-armv6l.tar.gz
else
echo "64 bit processor"
wget -nc https://go.dev/dl/$GOLANG_LATEST_STABLE_VERSION.linux-arm64.tar.gz $GOLANG_LATEST_STABLE_VERSION.linux-arm64.tar.gz
tar -zxvf /usr/local/$GOLANG_LATEST_STABLE_VERSION.linux-arm64.tar.gz
fi


echo export PATH=$PATH:/usr/local/go/bin >>  ~/.bashrc
echo export GOPATH=/home/talkkonnect/gocode >>  ~/.bashrc
echo export GOBIN=/home/talkkonnect/bin >>  ~/.bashrc
echo export GO111MODULE="auto" >>  ~/.bashrc
echo "alias tk='cd /home/talkkonnect/gocode/src/github.com/talkkonnect/iotalerter/'" >>  ~/.bashrc


## Set up GOENVIRONMENT
export PATH=$PATH:/usr/local/go/bin
export GOPATH=/home/talkkonnect/gocode
export GOBIN=/home/talkkonnect/bin

## Get the latest source code of talkkonnect from githu.com
cd $GOPATH
mkdir -p /home/talkkonnect/gocode/src/github.com/talkkonnect
cd /home/talkkonnect/gocode/src/github.com/talkkonnect
git clone https://github.com/talkkonnect/iotalerter
cd /home/talkkonnect/gocode/src/github.com/talkkonnect/iotalerter
go mod init
go mod tidy

## Build talkkonnect as binary
cd $GOPATH/src/github.com/talkkonnect/iotalerter
/usr/local/go/bin/go build -o /home/talkkonnect/bin/iotalerter cmd/main.go

exit


