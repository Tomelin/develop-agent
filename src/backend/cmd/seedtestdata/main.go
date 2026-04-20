package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

func main() {
	base := filepath.Join("testdata")
	files := []string{"projects.json", "billing_records.json"}

	for _, name := range files {
		path := filepath.Join(base, name)
		raw, err := os.ReadFile(path)
		if err != nil {
			panic(err)
		}
		var anyJSON any
		if err := json.Unmarshal(raw, &anyJSON); err != nil {
			panic(err)
		}
		fmt.Printf("fixture validated: %s\n", path)
	}
	fmt.Println("seed-test-data completed (fixtures ready for integration tests)")
}
