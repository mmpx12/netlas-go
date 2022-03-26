package cmd

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/spf13/cobra"
)

var raw bool
var faviconCmd = &cobra.Command{
	Use:   "favicon",
	Short: "Search from favicon",
	Long:  "Search from favicon",
	Run: func(cmd *cobra.Command, args []string) {
		url, _ := cmd.Flags().GetString("url")
		path, _ := cmd.Flags().GetString("path")
		apikey, _ := cmd.Flags().GetString("api-key")
		format, _ := cmd.Flags().GetString("format")
		if path == "" && url == "" {
			println("Error:\nMissing url or file\n")
			cmd.Usage()
			return
		}

		searchFavicon(url, path, apikey, format, colors, raw)
	},
}

func init() {
	rootCmd.AddCommand(faviconCmd)
	faviconCmd.Flags().StringP("url", "u", "", "Url of the favicon")
	faviconCmd.Flags().StringP("file", "F", "", "Path of the favicon file")
	faviconCmd.Flags().StringP("api-key", "a", "", "Specify api key (overwrite config file)")
	faviconCmd.Flags().StringP("format", "f", "yaml", "Select output format (yaml/json)")
	faviconCmd.Flags().BoolVarP(&colors, "no-colors", "n", false, "Disable colors")
	faviconCmd.Flags().BoolVarP(&raw, "raw-result", "r", false, "Print raw results")
}

func urlSum(url string) string {
	resp, err := http.Get(url)
	if err != nil {
		log.Fatalln(err)
	}
	defer resp.Body.Close()

	h := sha256.New()
	if _, err := io.Copy(h, resp.Body); err != nil {
		log.Fatal(err)
	}
	sum := fmt.Sprintf("%x", h.Sum(nil))
	return sum

}

func fileSum(path string) string {
	f, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		log.Fatal(err)
	}

	sum := fmt.Sprintf("%x", h.Sum(nil))
	return sum
}

func parse(data []byte, format string, colors bool) {
	query := `[.items[].data  | { "date": .scan_date, "ip": .ip, "domain": .domain}]|unique|.[]`
	result := Gojq(query, data)
	if format == "yaml" {
		var s []byte
		s, _ = json.Marshal(result)
		YamlPrint(s, colors)
	} else {
		JsonPrint(result, colors)
	}
}

func searchFavicon(url, path, apikey, format string, colors, raw bool) {
	var sum string
	if url != "" {
		sum = urlSum(url)
	} else if path != "" {
		sum = fileSum(path)
	}
	query := "https://app.netlas.io/api/responses/?q=favicon.hash_sha256:'" + string(sum) + "'"
	res := GetBody("GET", query, "", apikey)
	if format == "json" && raw {
		raw := json.RawMessage(string(res))
		JsonPrint(raw, colors)
	} else if format == "yaml" && raw {
		YamlPrint(res, colors)
	} else {
		parse(res, format, colors)
	}
}
