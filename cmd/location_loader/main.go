package main

import (
	"context"
	"fmt"
	"log"
	"os"

	places "cloud.google.com/go/maps/places/apiv1"
	"github.com/WesleyT4N/what-to-eat/internal/locationhistory"
	"github.com/joho/godotenv"
	"google.golang.org/api/option"
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

func NewGoogleMapsClient() (*places.Client, error) {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Load the API key from the environment
	apiKey := os.Getenv("GOOGLE_MAPS_API_KEY")

	if err != nil {
		log.Fatalf("fatal error: %s", err)
	}
	return places.NewClient(context.Background(), option.WithAPIKey(apiKey))
}

func main() {
	filename := parseArgs()
	fmt.Println("Loading Locations from", filename)

	locations, err := locationhistory.LoadLocationsFromDir(filename)
	if err != nil {
		panic(err)
	}
	fmt.Println("Loaded", len(locations), "locations")
	/* c, err := NewGoogleMapsClient() */
	/* if err != nil { */
	/* 	panic(err) */
	/* } */
	/* restaurants := locationhistory.GetRestaurants(locations, c) */
	/* fmt.Println("Found", len(restaurants), "restaurants") */
}
