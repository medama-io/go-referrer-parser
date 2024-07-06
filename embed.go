package referrer

import (
	_ "embed"
)

//go:embed data/referers.csv
var referrersCSV string
