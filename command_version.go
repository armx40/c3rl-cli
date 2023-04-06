package main

import (
	"fmt"

	"github.com/urfave/cli/v2"
)

const (
	command_version_major = 0
	command_version_minor = 1
)

func command_version_action(cCtx *cli.Context) (err error) {
	fmt.Printf("Version: %d.%d\n", command_version_major, command_version_minor)
	return
}
