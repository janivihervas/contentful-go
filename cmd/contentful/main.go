package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"net/url"
	"os"
	"strings"

	"github.com/janivihervas/contentful-go"
)

var (
	token   string
	spaceID string
	preview bool
)

func init() {
	flag.StringVar(&token, "token", "", "Contentful access token")
	flag.StringVar(&spaceID, "space", "", "Contentful space id")
	flag.BoolVar(&preview, "preview", false, "Whether to use the preview API or not")
}

func main() {
	flag.Parse()

	if token == "" || spaceID == "" {
		flag.Usage()
		os.Exit(1)
	}

	if len(flag.Args()) == 0 {
		fmt.Println("Must specify a query as a list of key=value pairs")
		os.Exit(1)
	}

	parameters := url.Values{}
	for _, arg := range flag.Args() {
		parts := strings.Split(arg, "=")
		if len(parts) != 2 {
			fmt.Println("Could not parse query:", arg)
			os.Exit(1)
		}
		parameters.Add(parts[0], parts[1])
	}

	cms := contentful.New(token, spaceID, preview)

	result := make(map[string]interface{})
	err := cms.Search(context.Background(), parameters, &result)
	if err != nil {
		fmt.Println("Client returned an error:", err)
		os.Exit(1)
	}

	bytes, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		fmt.Println("Could not parse result:", err)
		os.Exit(1)
	}

	fmt.Println(string(bytes))
}
