// Package format provides utilities for formating Strava data
package format

import (
	"encoding/json"
	"os"
)

func PrintAsJSON(data any) error {
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	return encoder.Encode(data)
}
