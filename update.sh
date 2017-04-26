#!/bin/sh
GIT_SSL_NO_VERIFY=1 git pull origin master && make.go -name sshd
