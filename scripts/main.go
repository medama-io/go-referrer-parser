package main

import (
	"encoding/csv"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
)

const (
	// Referer-parser dataset URL
	SOURCE_URL = "https://s3-eu-west-1.amazonaws.com/snowplow-hosted-assets/third-party/referer-parser/referers-latest.json"

	// Output CSV file path
	CSV_FILE_PATH = "./data/referers.csv"
)

var data map[string]interface{} = make(map[string]interface{})

// Fetch dataset from https://github.com/snowplow-referer-parser/referer-parser
func getData() {
	resp, err := http.Get(SOURCE_URL)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	err = json.Unmarshal(body, &data)
	if err != nil {
		log.Fatal(err)
	}
}

func convertToCSV() {
	f, err := os.Create(CSV_FILE_PATH)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	d := csv.NewWriter(f)
	defer d.Flush()

	// e.g. unknown, search, social, email, paid
	for group, groupValue := range data {
		// e.g. Google, Yandex Maps, Yahoo!
		for name, nameValue := range groupValue.(map[string]interface{}) {
			for _, domain := range nameValue.(map[string]interface{})["domains"].([]interface{}) {
				// Write a new row to referers.csv.
				d.Write([]string{group, name, domain.(string)})
			}
		}
	}
}

func main() {
	getData()
	convertToCSV()
}
