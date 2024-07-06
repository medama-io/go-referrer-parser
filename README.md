# go-referrer-parser

Convert HTTP Referer URLs into grouped attributions.

## Installation

```bash
go get github.com/medama-io/go-referrer-parser
```

## Usage

```go
package main

import (
	"fmt"
	"github.com/medama-io/go-referrer-parser"
)

func main() {
	// Create a new parser on start up.
	parser := referrer.NewParser()

	// Parse a referrer URL.
	name, err := parser.Parse("https://www.google.co.uk")
	
	fmt.Println(name) // "Google"
}
