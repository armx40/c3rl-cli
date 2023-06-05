package main

import (
	"fmt"

	"github.com/AlecAivazis/survey/v2"
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

func command_auth_verify(cCtx *cli.Context) error {

	err := command_auth_functions_echo()
	if err != nil {
		return fmt.Errorf("failed to verify auth status")
	}

	fmt.Printf("auth status verified\n")

	return nil
}
