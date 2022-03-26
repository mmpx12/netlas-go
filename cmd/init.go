package cmd

import (
	"fmt"
	"os"

	"github.com/AlecAivazis/survey/v2"
	"github.com/AlecAivazis/survey/v2/terminal"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Save api key and server to file",
	Long:  "Save api key and server to file",
	Run: func(cmd *cobra.Command, args []string) {
		CreateConfig()
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}

func CreateConfig() {
	type Answers struct {
		Server string
		Token  string
	}
	answer := Answers{}
	var question = []*survey.Question{
		{
			Name: "server",
			Prompt: &survey.Input{
				Message: "Netlas server (domain or ip)",
				Help:    "For non standard http(s) port use <domain/ip:port>",
				Default: "https://app.netlas.io",
			},
		},
		{
			Name: "token",
			Prompt: &survey.Password{
				Message: "Netlas api key: ",
			},
		},
	}
	err := survey.Ask(question, &answer, survey.WithValidator(survey.Required))
	if err == terminal.InterruptErr {
		fmt.Println("interrupted")
		os.Exit(1)
	}

	viper.Set("host", answer.Server)
	viper.Set("token", answer.Token)
	viper.WriteConfig()
	fmt.Println("\033[1;32m✓ \033[0;0m Configuration updated")
}
