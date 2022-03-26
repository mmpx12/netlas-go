package cmd

import (
	"encoding/json"
	"strconv"

	"github.com/spf13/cobra"
)

var domainCmd = &cobra.Command{
	Use:   "domain",
	Short: "Search for domain",
	Long:  "Search for domain",
	Run: func(cmd *cobra.Command, args []string) {
		apikey, _ := cmd.Flags().GetString("api-key")
		format, _ := cmd.Flags().GetString("format")
		limit, _ := cmd.Flags().GetInt("page")
		domain, _ := cmd.Flags().GetString("domain")
		if domain == "" {
			println("Error:\nMissing -d or --domain flag\n")
			cmd.Usage()
			return
		}
		GetDomain(domain, limit, format, apikey, raw)
	},
}

func init() {
	rootCmd.AddCommand(domainCmd)
	domainCmd.Flags().StringP("api-key", "a", "", "Specify api key (overwrite config file)")
	domainCmd.Flags().StringP("format", "f", "json", "Select output format (yaml/json)")
	domainCmd.Flags().StringP("domain", "d", "", "Serch domain")
	domainCmd.Flags().BoolVarP(&colors, "no-colors", "n", false, "Disable colors")
	domainCmd.Flags().IntP("page", "p", 0, "Page number (default 1)")
	domainCmd.Flags().BoolVarP(&raw, "raw-result", "r", false, "Print raw results")
}

func GetDomain(domain string, limit int, format, apikey string, raw bool) {
	var query string
	if limit == 0 {
		query = "https://app.netlas.io/api/domains/?q=" + domain + "&start=" + strconv.Itoa(limit)
	} else {
		limit = limit * 20
		query = "https://app.netlas.io/api/domains/?q=" + domain + "&start=" + strconv.Itoa(limit)
	}
	res := GetBody("GET", query, "", apikey)
	if format == "json" && raw {
		raw := json.RawMessage(string(res))
		JsonPrint(raw, colors)
	} else if format == "yaml" && raw {
		YamlPrint(res, colors)
	} else {
		jqquery := `[.items[]| {"ip": (try (.data.a|.[])// (try (.data.txt|.[]) // "Unknow")), "domain": .data.domain}]|reduce .[] as $d (null; .[$d.ip] += [$d.domain])`
		Parse(res, format, jqquery, colors)
	}
}
