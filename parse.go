package referrer

import (
	"encoding/csv"
	"strings"
)

type Parser struct {
	// referers is a map of referer groups to a map of referer names to a slice of domains.
	referrers map[string]string

	// spammers is a map of spammer domains to their names.
	spammers map[string]struct{}
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

	// Populate spammers map with spammer domains.
	spammers := make(map[string]struct{})
	spammerLines := strings.Split(spammersList, "\n")
	for _, line := range spammerLines {
		line = strings.TrimSpace(line)
		if line != "" {
			spammers[line] = struct{}{}
		}
	}

	return &Parser{
		referrers: referrers,
		spammers:  spammers,
	}, nil
}

// Parse returns the referer name for a given referer domain. If the domain is not found, it returns nil.
func (p *Parser) Parse(domain string) string {
	domain = strings.ToLower(strings.TrimSpace(domain))

	// Check for exact match
	if name, ok := p.referrers[domain]; ok {
		return name
	}

	// Check for a match stripping the leading "www."
	if strings.HasPrefix(domain, "www.") {
		if name, ok := p.referrers[domain[4:]]; ok {
			return name
		}
	}

	return ""
}

// IsSpam checks if the given domain is in the spammer list.
func (p *Parser) IsSpam(domain string) bool {
	domain = strings.ToLower(strings.TrimSpace(domain))
	_, exists := p.spammers[domain]
	if exists {
		return true
	}

	// Check for a match stripping the leading "www."
	if strings.HasPrefix(domain, "www.") {
		_, exists = p.spammers[domain[4:]]
		return exists
	}

	return false
}
