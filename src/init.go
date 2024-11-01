package challenge2016

import (
	"fmt"
)

// NewDistributor initializes a new Distributor with the updated Permission structure
func NewDistributor(name string, parent *Distributor) *Distributor {
	return &Distributor{
		Name:   name,
		Parent: parent,
		Permissions: Permission{
			Include: make(map[string]map[string]map[string]bool), // Updated structure for Include
			Exclude: make(map[string]map[string]map[string]bool), // Updated structure for Exclude
		},
	}
}

func (d *Distributor) AddInclude(locationObj LocationObject) {
	// Ensure country and province maps exist
	if _, ok := d.Permissions.Include[locationObj.Country]; !ok {
		d.Permissions.Include[locationObj.Country] = make(map[string]map[string]bool)
	}
	if _, ok := d.Permissions.Include[locationObj.Country][locationObj.Province]; !ok {
		d.Permissions.Include[locationObj.Country][locationObj.Province] = make(map[string]bool)
	}

	// Add city to the province if not already included
	if _, cityExists := d.Permissions.Include[locationObj.Country][locationObj.Province][locationObj.City]; cityExists {
		fmt.Printf("City %s in province %s of country %s already included.\n", locationObj.City, locationObj.Province, locationObj.Country)
	} else {
		d.Permissions.Include[locationObj.Country][locationObj.Province][locationObj.City] = true
	}
}

func (d *Distributor) AddExclude(locationObj LocationObject) {
	// Ensure country and province maps exist
	if _, ok := d.Permissions.Exclude[locationObj.Country]; !ok {
		d.Permissions.Exclude[locationObj.Country] = make(map[string]map[string]bool)
	}
	if _, ok := d.Permissions.Exclude[locationObj.Country][locationObj.Province]; !ok {
		d.Permissions.Exclude[locationObj.Country][locationObj.Province] = make(map[string]bool)
	}

	// Add city to the province if not already excluded
	if _, cityExists := d.Permissions.Exclude[locationObj.Country][locationObj.Province][locationObj.City]; cityExists {
		fmt.Printf("City %s in province %s of country %s already excluded.\n", locationObj.City, locationObj.Province, locationObj.Country)
	} else {
		d.Permissions.Exclude[locationObj.Country][locationObj.Province][locationObj.City] = true
	}
}

func (d *Distributor) CanDistribute(locationObj LocationObject) bool {
	// Check if the location is explicitly excluded
	if provinces, ok := d.Permissions.Exclude[locationObj.Country]; ok {
		if cities, exists := provinces[locationObj.Province]; exists {
			if _, cityExists := cities[locationObj.City]; cityExists {
				return false
			}
		}
	}

	// Check if the country is included
	if provinces, ok := d.Permissions.Include[locationObj.Country]; ok {
		// If no specific province or city is checked, return true
		if locationObj.Province == "" && locationObj.City == "" {
			return true
		}

		if locationObj.Province != "" {
			if cities, exists := provinces[locationObj.Province]; exists {
				// If the province is included, check for city
				if locationObj.City == "" || cities[locationObj.City] {
					return true
				}
			} else {
				return true // country is included and province not, also neither country not province in exclusion so distributor can distribute to this province
			}
		} else {
			// If only the country is included and no specific province is checked, allow distribution
			return true
		}
	}

	return false // If none of the conditions are met, return false
}

func (d *Distributor) Includes(locationObj LocationObject) bool {
	// Check if the location is included
	if provinces, ok := d.Permissions.Include[locationObj.Country]; ok {
		if cities, exists := provinces[locationObj.Province]; exists {
			if _, cityExists := cities[locationObj.City]; cityExists {
				return true
			}
		}
	}
	return false
}

func (d *Distributor) Excludes(locationObj LocationObject) bool {
	// Check if the location is excluded
	if provinces, ok := d.Permissions.Exclude[locationObj.Country]; ok {
		if cities, exists := provinces[locationObj.Province]; exists {
			if _, cityExists := cities[locationObj.City]; cityExists {
				return true
			}
		}
	}
	return false
}

func (d *Distributor) RemoveLocation(locationObj LocationObject, targetPermissionType string) {
	// Choose either the Include or Exclude map directly
	var targetMap map[string]map[string]map[string]bool
	if targetPermissionType == "EXCLUDE" {
		targetMap = d.Permissions.Exclude
	} else {
		targetMap = d.Permissions.Include
	}
	removeFromMap(targetMap, locationObj)
}

// Helper function to remove the location from the given map
func removeFromMap(targetMap map[string]map[string]map[string]bool, locationObj LocationObject) {
	// Check if the country exists
	if provinces, ok := targetMap[locationObj.Country]; ok {
		// Check if the province exists in the country
		if cities, exists := provinces[locationObj.Province]; exists {
			// Remove the city if it exists
			delete(cities, locationObj.City)
			// If no cities remain, remove the province
			if len(cities) == 0 {
				delete(provinces, locationObj.Province)
			}
			// If no provinces remain, remove the country
			if len(provinces) == 0 {
				delete(targetMap, locationObj.Country)
			}
		}
	}
}
