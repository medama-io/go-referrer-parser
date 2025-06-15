package main

import (
	"encoding/csv"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"slices"
	"strings"
)

const (
	// Referer-parser dataset URL
	SOURCE_URL = "https://s3-eu-west-1.amazonaws.com/snowplow-hosted-assets/third-party/referer-parser/referers-latest.json"

	// Output CSV file path
	CSV_FILE_PATH = "./data/referers.csv"
)

var groupData map[string]any = make(map[string]any)

// Fetch dataset from https://github.com/snowplow-referer-parser/referer-parser
func getGroupData() {
	resp, err := http.Get(SOURCE_URL)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	err = json.Unmarshal(body, &groupData)
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

	var records [][]string

	// e.g. unknown, search, social, email, paid
	for group, groupValue := range groupData {
		// e.g. Google, Yandex Maps, Yahoo!
		for name, nameValue := range groupValue.(map[string]any) {
			for _, domain := range nameValue.(map[string]any)["domains"].([]any) {
				// Write a new row to referers.csv.
				records = append(records, []string{group, name, domain.(string)})
			}
		}
	}

	// Sort records
	slices.SortFunc(records, func(a, b []string) int {
		// Sort by group, then name, then domain
		if a[0] == b[0] {
			if a[1] == b[1] {
				return strings.Compare(a[2], b[2])
			}
			return strings.Compare(a[1], b[1])
		}
		return strings.Compare(a[0], b[0])
	})

	// Write sorted records to CSV
	d := csv.NewWriter(f)
	defer d.Flush()

	for _, record := range records {
		if err := d.Write(record); err != nil {
			log.Fatal(err)
		}
	}

	log.Printf("CSV file written to %s\n", CSV_FILE_PATH)
}

// Fetch list from Matomo's referrer spam list.
func writeSpamList() {
	resp, err := http.Get("https://raw.githubusercontent.com/matomo-org/referrer-spam-list/master/spammers.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	file, err := os.Create("./data/spammers.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	_, err = file.Write(body)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Spam list written to ./data/spammers.txt")
}

func main() {
	getGroupData()
	convertToCSV()

	writeSpamList()
}
