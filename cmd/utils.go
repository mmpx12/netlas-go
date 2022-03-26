package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/fatih/color"
	"github.com/goccy/go-yaml"
	"github.com/goccy/go-yaml/lexer"
	"github.com/goccy/go-yaml/printer"
	"github.com/itchyny/gojq"
	"github.com/mattn/go-colorable"
	js "github.com/nwidger/jsoncolor"
	"github.com/spf13/viper"
)

func GetBody(method, url, param, apikey string) []byte {
	client := http.Client{}
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		panic(err)
	}
	var apiKey string
	if apikey == "" {
		apiKey = viper.GetString("token")
	} else {
		apiKey = apikey
	}
	req.Header = http.Header{
		"X-Api-Key":    []string{apiKey},
		"Content-Type": []string{"application/json"},
	}
	res, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()
	b, _ := io.ReadAll(res.Body)
	return b
}

func format(attr color.Attribute) string {
	return fmt.Sprintf("\x1b[%dm", attr)
}

func JsonPrint(data interface{}, colors bool) {
	var s []byte
	if colors {
		s, _ = json.MarshalIndent(data, "", "  ")

	} else {
		f := js.NewFormatter()
		s, _ = js.MarshalIndentWithFormatter(data, "", "  ", f)
	}
	fmt.Println(string(s))

}

func YamlPrint(data []byte, colorz bool) {
	x, err := yaml.JSONToYAML(data)
	if err != nil {
		panic(err)
	}
	tokens := lexer.Tokenize(string(x))
	var p printer.Printer
	if !colorz {
		p.Bool = func() *printer.Property {
			return &printer.Property{
				Prefix: format(color.FgHiMagenta),
				Suffix: format(color.Reset),
			}
		}
		p.Number = func() *printer.Property {
			return &printer.Property{
				Prefix: format(color.FgHiRed),
				Suffix: format(color.Reset),
			}
		}
		p.MapKey = func() *printer.Property {
			return &printer.Property{
				Prefix: format(color.FgHiCyan),
				Suffix: format(color.Reset),
			}
		}
		p.String = func() *printer.Property {
			return &printer.Property{
				Prefix: format(color.FgHiGreen),
				Suffix: format(color.Reset),
			}
		}
	}
	writer := colorable.NewColorableStdout()
	writer.Write([]byte(p.PrintTokens(tokens) + "\n"))
}

func Gojq(querystr string, data []byte) []interface{} {
	query, err := gojq.Parse(querystr)
	if err != nil {
		panic(err)
	}
	var input interface{}
	result := make([]interface{}, 0)
	json.Unmarshal(data, &input)
	value := query.Run(input)
	for {
		v, ok := value.Next()
		if !ok {
			break
		}
		if err, ok := v.(error); ok {
			panic(err)
		}
		result = append(result, v)
	}
	return result
}

func Parse(data []byte, format, query string, colors bool) {
	result := Gojq(query, data)
	if format == "yaml" {
		var s []byte
		s, _ = json.Marshal(result)
		YamlPrint(s, colors)
	} else {
		JsonPrint(result, colors)
	}
}
