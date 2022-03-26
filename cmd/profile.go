package cmd

import (
	"encoding/json"

	"github.com/spf13/cobra"
)

// profileCmd represents the profile command
var colors bool
var profileCmd = &cobra.Command{
	Use:   "profile",
	Short: "Get profile info",
	Long:  "Get profile info",
	Run: func(cmd *cobra.Command, args []string) {
		apikey, _ := cmd.Flags().GetString("api-key")
		format, _ := cmd.Flags().GetString("format")
		User(apikey, format, colors)
	},
}

func init() {
	rootCmd.AddCommand(profileCmd)
	profileCmd.Flags().StringP("api-key", "a", "", "Specify api key (overwrite config file)")
	profileCmd.Flags().StringP("format", "f", "yaml", "Select output format (yaml/json)")
	profileCmd.Flags().BoolVarP(&colors, "no-colors", "n", false, "Disable colors")
}

func User(apikey string, outputFormat string, colors bool) {
	url := "https://app.netlas.io/api/users/profile/"
	b := GetBody("GET", url, "", apikey)

	if outputFormat == "json" {
		raw := json.RawMessage(string(b))
		JsonPrint(raw, colors)
	} else {
		YamlPrint(b, colors)
	}
}
