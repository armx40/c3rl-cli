package main

import (
	"fmt"
	"strconv"
	"time"

	"github.com/urfave/cli/v2"
)

var (
	command_version_major        = "0"
	command_version_minor        = "0"
	command_version_patch        = "0"
	command_version_build_number = "0"
	command_version_time_unix    = "0"
)

func command_version_action(cCtx *cli.Context) (err error) {

	/* get time */
	i, err := strconv.ParseInt(command_version_time_unix, 10, 64)
	if err != nil {
		panic(err)
	}
	tm := time.Unix(i, 0)
	tm_formatted := tm.Format("02-01-2006")
	/**/

	fmt.Printf("Version: %s.%s.%s\nBuild: %s\n(%s)\n", command_version_major, command_version_minor, command_version_patch, command_version_build_number, tm_formatted)
	return
}
