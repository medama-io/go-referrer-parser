package referrer

import (
	"encoding/csv"
	"strings"
)

type Parser struct {
	// referers is a map of referer groups to a map of referer names to a slice of domains.
	referrers map[string]string
}

// NewParser creates a new Parser instance.
func NewParser() (*Parser, error) {
	referrers := make(map[string]string)

	// Populate referrers map with referer domains --> referrer names.
	r := csv.NewReader(strings.NewReader(referrersCSV))

	records, err := r.ReadAll()
	if err != nil {
		return nil, err
	}
	for _, row := range records {
		referrers[row[2]] = row[1]
	}

	return &Parser{
		referrers: referrers,
	}, nil
}

// Parse returns the referer name for a given referer domain. If the domain is not found, it returns nil.
func (p *Parser) Parse(domain string) string {
	domain = strings.ToLower(strings.TrimSpace(domain))

	// Check for exact match
	if name, ok := p.referrers[domain]; ok {
		return name
	}
	return ""
}
