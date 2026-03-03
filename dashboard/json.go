package dashboard

import (
	"encoding/json"
	"fmt"

	"github.com/ZadgeIsCool/gitpulse/analysis"
)

// PrintJSON outputs the report as formatted JSON.
func PrintJSON(r *analysis.Report) error {
	data, err := json.MarshalIndent(r, "", "  ")
	if err != nil {
		return err
	}
	fmt.Println(string(data))
	return nil
}
