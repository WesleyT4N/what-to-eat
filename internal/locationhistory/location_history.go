package locationhistory

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"googlemaps.github.io/maps"
)

type Location struct {
	LatitudeE7            int        `json:"latitudeE7"`
	LongitudeE7           int        `json:"longitudeE7"`
	PlaceID               string     `json:"placeId"`
	Address               string     `json:"address"`
	Name                  string     `json:"name"`
	SemanticType          string     `json:"semanticType"`
	SourceInfo            SourceInfo `json:"sourceInfo"`
	LocationConfidence    float64    `json:"locationConfidence"`
	CalibratedProbability float64    `json:"calibratedProbability"`
}

type LocationHistory struct {
	Locations []Location `json:"timelineObjects"`
}

type SourceInfo struct {
	DeviceTag int `json:"deviceTag"`
}

type LocationFrequency struct {
	Location  Location
	Frequency int
}

func LoadLocationsFromDir(dir string) ([]LocationFrequency, error) {
	locations := map[string]LocationFrequency{}
	filepath.WalkDir(dir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if strings.HasSuffix(d.Name(), "_filtered.json") {
			fileLocations, err := LoadLocationsFromFile(path)
			if err != nil {
				return err
			}

			for _, location := range fileLocations {
				if location.LocationConfidence > 50 {
					if _, ok := locations[location.PlaceID]; ok {
						locations[location.PlaceID] = LocationFrequency{
							Location:  location,
							Frequency: locations[location.PlaceID].Frequency + 1,
						}
					} else {
						locations[location.PlaceID] = LocationFrequency{
							Location:  location,
							Frequency: 1,
						}
					}
				}
			}
		}
		return nil
	})

	var uniqueLocations []LocationFrequency
	for _, location := range locations {
		uniqueLocations = append(uniqueLocations, location)
	}

	return uniqueLocations, nil
}

func LoadLocationsFromFile(filename string) ([]Location, error) {
	jsonFile, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var result LocationHistory
	json.Unmarshal(jsonFile, &result)

	return result.Locations, nil
}

type Restaurant struct {
	Name    string
	Address string
	PlaceID string
	Visits  int
}

var restaurantPlaceTypes map[string]bool = map[string]bool{
	"restaurant":    true,
	"food":          true,
	"meal_takeaway": true,
	"meal_delivery": true,
}

func isRestaurant(place *maps.PlaceDetailsResult) bool {
	for _, placeType := range place.Types {
		if _, ok := restaurantPlaceTypes[placeType]; ok {
			return true
		}
	}
	return false
}

func GetRestaurants(locations []LocationFrequency, c *maps.Client) []Restaurant {
	var restaurants []Restaurant
	for _, location := range locations {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		place, err := c.PlaceDetails(ctx, &maps.PlaceDetailsRequest{
			PlaceID: location.Location.PlaceID,
		})

		if err != nil {
			log.Printf("Error getting place details for %s: %v", location.Location.PlaceID, err)
			continue
		}

		if isRestaurant(&place) {
			restaurants = append(restaurants, Restaurant{
				Name:    place.Name,
				Address: place.FormattedAddress,
				PlaceID: place.PlaceID,
				Visits:  location.Frequency,
			})
			fmt.Printf("Found restaurant %s at %s with %d visits\n", place.Name, place.FormattedAddress, location.Frequency)
		}
	}
	return restaurants
}
