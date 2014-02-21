#!/bin/bash

rsync --verbose --progress --recursive --update --exclude site --exclude .git --exclude .DS_Store . direct.txgruppi.com:go/src/github.com/txgruppi/site
ssh direct.txgruppi.com 'cd /home/txgruppi/go/src/github.com/txgruppi/site/; if [ ./site -ot ./main.go ]; then GOPATH=/home/txgruppi/go/ go build -ldflags "-s" -v; fi; sudo ./restart.sh'
