#!/bin/bash
git pull
go build -o c3rl-cli
strip ./c3rl-cli
sudo install -o $(id -nu) -g $(id -ng) -m 0755 c3rl-cli /usr/local/bin/c3rl-cli 