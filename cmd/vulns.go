package cmd

import (
	"encoding/json"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
)

var vulnsCmd = &cobra.Command{
	Use:   "vulns",
	Short: "Search host by cve",
	Long:  "Search host by cve",
	Run: func(cmd *cobra.Command, args []string) {
		apikey, _ := cmd.Flags().GetString("api-key")
		format, _ := cmd.Flags().GetString("format")
		limit, _ := cmd.Flags().GetInt("page")
		id, _ := cmd.Flags().GetString("cve-id")
		exploit, _ := cmd.Flags().GetBool("has-exploit")
		score, _ := cmd.Flags().GetString("min-score")
		if id == "" && score == "" && !exploit {
			println("Error:\n")
			cmd.Usage()
			return
		}

		GetVulns(id, score, limit, format, apikey, raw, exploit)
	},
}

func init() {
	rootCmd.AddCommand(vulnsCmd)
	vulnsCmd.Flags().StringP("api-key", "a", "", "Specify api key (overwrite config file)")
	vulnsCmd.Flags().StringP("format", "f", "json", "Select output format (yaml/json)")
	vulnsCmd.Flags().StringP("cve-id", "i", "", "Serch for cve id")
	vulnsCmd.Flags().BoolP("has-exploit", "x", false, "Serch host with cve that have exploit")
	vulnsCmd.Flags().StringP("min-score", "m", "0", "Minium score of cve to search")
	vulnsCmd.Flags().BoolVarP(&colors, "no-colors", "n", false, "Disable colors")
	vulnsCmd.Flags().IntP("page", "p", 0, "Page number (default 1)")
	vulnsCmd.Flags().BoolVarP(&raw, "raw-result", "r", false, "Print raw results")
}

func GetVulns(id, score string, limit int, format, apikey string, raw, exploit bool) {
	var query, jqquery strings.Builder
	if id != "" {
		query.WriteString("https://app.netlas.io/api/responses/?q=cve.name:" + id)
		jqquery.WriteString(`{"` + id + `":([.items[].data.ip]|unique)}`)
	} else if score != "0" {
		query.WriteString("https://app.netlas.io/api/responses/?q=(cve.base_score:>" + score + ")")
		jqquery.WriteString(`[.items[].data| {"ip": .ip, "cve": [ (.cve[]|if (.base_score|tonumber) > ` + score + ` then (.name + " ("+.base_score+ ")") else empty end)]}|.]|unique[]`)
	}
	if exploit {
		if query.Len() == 0 {
			query.WriteString("https://app.netlas.io/api/responses/?q=cve.has_exploit:true")
		} else {
			query.WriteString("AND(cve.has_exploit:true)")
			jqquery.Reset()
			jqquery.WriteString(`[.items[].data| {"ip": .ip, "cve": [ (.cve[]|if (.base_score|tonumber) > ` + score + ` and .has_exploit == true then (.name + " ("+.base_score+ ")") else empty end)]}|.]|unique[]`)
		}
	}
	if limit != 0 {
		limit = limit * 20
		query.WriteString("&start=" + strconv.Itoa(limit))
	}
	res := GetBody("GET", query.String(), "", apikey)
	if format == "json" && raw {
		raw := json.RawMessage(string(res))
		JsonPrint(raw, colors)
	} else if format == "yaml" && raw {
		YamlPrint(res, colors)
	} else {
		Parse(res, format, jqquery.String(), colors)
	}
}
