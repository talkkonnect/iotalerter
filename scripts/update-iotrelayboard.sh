#!/bin/bash

SERVICE="iotalerter"
BACKUPXML=iotalerter-$(date +"%Y%m%d-%H%M%S").xml

if pgrep -x "$SERVICE" >/dev/null
then
    echo "$SERVICE is running"
    systemctl stop iotalerter
else
    echo "$SERVICE stopped"
fi

if [[ -f "/home/talkkonnect/bin/iotalerter" ]]
then
	echo "removing /home/talkkonnect/bin/iotalerter binary"
	rm /home/talkkonnect/bin/iotalerter
fi

if [[ -f "/home/talkkonnect/gocode/src/github.com/talkkonnect/iotalerter/iotalerter.xml" ]]
then
	echo "copying iotalerter.xml for safe keeping to /root/"$BACKUPXML
	cp /home/talkkonnect/gocode/src/github.com/talkkonnect/iotalerter/iotalerter.xml /root/$BACKUPXML
fi

rm -rf /home/talkkonnect/gocode/src/github.old
rm -rf /home/talkkonnect/gocode/src/google.golang.org
rm -rf /home/talkkonnect/gocode/src/golang.org
rm -rf  /home/talkkonnect/gocode/src/github.com
rm -rf  /home/talkkonnect/bin/iotalerter


## Create the necessary directoy structure under /home/talkkonnect/
mkdir -p /home/talkkonnect/gocode
mkdir -p /home/talkkonnect/gocode/src
mkdir -p /home/talkkonnect/gocode/src/github.com


## Added this block to update to the latest version of golang so the update doesnt break talkkonnect
rm -rf /usr/local/go
cd /usr/local
cd /usr/local

## Check Latest of GOLANG 64 Bit Version for Raspberry Pi
GOLANG_LATEST_STABLE_VERSION=$(curl -s https://go.dev/VERSION?m=text | grep go)
cputype=`lscpu | grep Architecture | cut -d ":" -f 2 | sed 's/ //g'`
bitsize=`getconf LONG_BIT`

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

## Set up GOENVIRONMENT
export PATH=$PATH:/usr/local/go/bin
export GOPATH=/home/talkkonnect/gocode
export GOBIN=/home/talkkonnect/bin
export GO111MODULE="auto"

## Get the latest source code of talkkonnect from github.com
echo "getting iotalerter with go get"
cd $GOPATH
go get -v github.com/talkkonnect/iotalerter

## Build talkkonnect as binary
cd $GOPATH/src/github.com/talkkonnect/iotalerter
/usr/local/go/bin/go build -o /home/talkkonnect/bin/iotalerter /home/talkkonnect/gocode/src/github.com/talkkonnect/iotalerter/cmd/main.go

cp /root/$BACKUPXML /home/talkkonnect/gocode/src/github.com/talkkonnect/iotalerter/iotalerter.xml

if pgrep -x "$SERVICE" >/dev/null
then
    echo "$SERVICE is running I will stop it please start talkkonnect manually"
    systemctl stop iotalerter
else
    echo "$SERVICE is stopped now restarting iotalerter"
    systemctl start iotalerter
fi

## Notify User
echo "=> Finished Updating iotalerter"
echo "=> Updated iotalerter binary is in /home/talkkonect/bin"

exit
