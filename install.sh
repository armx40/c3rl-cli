#!/bin/bash
go build -o c3rl-cli
strip ./c3rl-cli
sudo install -o root -g root -m 0755 c3rl-cli /usr/local/bin/c3rl-cli 