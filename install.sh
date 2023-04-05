#!/bin/bash
go build -o c3rl-cli
strip ./c3rl-cli
sudo ln -s $(pwd)/c3rl-cli /usr/bin/c3rl-cli 