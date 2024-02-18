package main

import (
	"fmt"
	"os"

	"github.com/WesleyT4N/what-to-eat/internal/maps"
)

func parseArgs() string {
	// parse filepath from first arg
	if len(os.Args) < 2 {
		fmt.Println("Usage: location_loader <dir_path>")
		os.Exit(1)
	}
	filePath := os.Args[1]
	return filePath
}

func main() {
	filename := parseArgs()
	fmt.Println("Loading Locations from", filename)

	locations, err := maps.LoadLocationsFromDir(filename)
	if err != nil {
		panic(err)
	}
	for _, location := range locations {
		loc := location.Location
		freq := location.Frequency
		fmt.Println("Location", loc.Name, "at", loc.Address, "with confidence", loc.LocationConfidence, ". Appeared", freq, "times")
	}

	fmt.Println("Loaded", len(locations), "locations")
}
