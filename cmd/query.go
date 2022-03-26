package cmd

import (
	"encoding/json"
	"strconv"

	"github.com/spf13/cobra"
)

var queryCmd = &cobra.Command{
	Use:   "query",
	Short: "Search for custom query",
	Long:  "Search for custom query",
	Run: func(cmd *cobra.Command, args []string) {
		apikey, _ := cmd.Flags().GetString("api-key")
		format, _ := cmd.Flags().GetString("format")
		limit, _ := cmd.Flags().GetInt("page")
		query, _ := cmd.Flags().GetString("query")
		if query == "" {
			println("Error:\nMissing -q or --query flag\n")
			cmd.Usage()
			return
		}

		GetQuery(query, limit, format, apikey, raw)
	},
}

func init() {
	rootCmd.AddCommand(queryCmd)
	queryCmd.Flags().StringP("api-key", "a", "", "Specify api key (overwrite config file)")
	queryCmd.Flags().StringP("format", "f", "json", "Select output format (yaml/json)")
	queryCmd.Flags().StringP("query", "d", "", "Serch domain")
	queryCmd.Flags().BoolVarP(&colors, "no-colors", "n", false, "Disable colors")
	queryCmd.Flags().IntP("page", "p", 0, "Page number (default 1)")
	queryCmd.Flags().BoolVarP(&raw, "raw-result", "r", false, "Print raw results")
}

func GetQuery(query string, limit int, format, apikey string, raw bool) {
	var querystr string
	if limit == 0 {
		querystr = "https://app.netlas.io/api/domains/?q=" + query + "&start=" + strconv.Itoa(limit)
	} else {
		limit = limit * 20
		querystr = "https://app.netlas.io/api/domains/?q=" + query + "&start=" + strconv.Itoa(limit)
	}
	res := GetBody("GET", querystr, "", apikey)
	if format == "json" {
		raw := json.RawMessage(string(res))
		JsonPrint(raw, colors)
	} else {
		YamlPrint(res, colors)
	}
}
