package main

import (
	challenge2016 "challenge2016/src"
	"fmt"
	"os"
)

func main() {
	// Load regions from CSV
	err := challenge2016.LoadRegions("cities.csv")
	if err != nil {
		fmt.Println("Error loading regions:", err)
		os.Exit(1)
	}

	// challenge2016.DumpMapToFile()
	// Run the interactive CLI
	challenge2016.RunInteractive()
}
