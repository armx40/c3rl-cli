package main

import (
	"fmt"

	"github.com/AlecAivazis/survey/v2"
	"github.com/fatih/color"
	"github.com/urfave/cli/v2"
)

var command_auth_login_username string
var command_auth_login_password string

func command_auth_subcommands() (commands cli.Commands) {

	commands = cli.Commands{{

		Name:    "login",
		Aliases: []string{"l"},
		Usage:   "Log in to your account",
		Action:  command_auth_login,
	}, {

		Name:    "verify",
		Aliases: []string{"v"},
		Usage:   "Verify the auth status",
		Action:  command_auth_verify,
	},
	}

	return commands
}

func command_auth_login(cCtx *cli.Context) error {
	var qs = []*survey.Question{
		{
			Name: "Username",
			Prompt: &survey.Input{
				Message: "Username:",
				Default: "",
			},
		},
		{
			Name: "Password",
			Prompt: &survey.Password{
				Message: "Password:",
			},
		},
	}

	var data AuthLoginSurveyAnswerPayload

	err := survey.Ask(qs, &data)
	if err != nil {
		fmt.Println(err)
		return err
	}

	err = command_auth_functions_login(data.Username, data.Password)
	if err != nil {
		return err
	}

	return nil
}

func command_auth_verify(cCtx *cli.Context) (err error) {

	err = command_auth_functions_echo()
	if err != nil {
		print_text := color.New(color.FgWhite)
		print_text.Printf("Status: ")
		print_text.Add(color.FgRed)
		print_text.Printf("FAIL\n")

		return fmt.Errorf("failed to verify auth status")
	}

	/* get auth data */
	// auth_data, err := command_auth_functions_get_auth_data()
	// if err != nil {
	// 	return
	// }
	/* */

	print_text := color.New(color.FgWhite)
	print_text.Printf("Status: ")
	print_text.Add(color.FgGreen)
	print_text.Printf("OK\n")

	return nil
}
