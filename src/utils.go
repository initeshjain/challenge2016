package challenge2016

import (
	"bufio"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
)

// LoadRegions loads regions from a CSV file and populates the Map structure.
func LoadRegions(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	_, err = reader.Read() // Skip the header row
	if err != nil {
		return err
	}

	// Initialize the main map
	Map = make(map[string]map[string]map[string]bool)

	for {
		record, err := reader.Read()
		if err != nil {
			if err == io.EOF {
				break // Break the loop on EOF
			}
			return err
		}

		// Extract relevant fields from the record
		country := UpperCaseAndTrimSpace(record[5])  // Country Name
		province := UpperCaseAndTrimSpace(record[4]) // Province Name
		city := UpperCaseAndTrimSpace(record[3])     // City Name

		// Initialize nested maps if they don't exist
		if _, countryExists := Map[country]; !countryExists {
			Map[country] = make(map[string]map[string]bool)
		}
		if _, provinceExists := Map[country][province]; !provinceExists {
			Map[country][province] = make(map[string]bool)
		}

		// Add the city to the province map
		Map[country][province][city] = true
	}

	return nil
}

// IsValidLocation checks if the given location is valid based on Map.
func IsValidLocation(location string) bool {
	location = UpperCaseAndTrimSpace(location)

	// Split the location into components
	parts := strings.Split(location, ",")
	if len(parts) < 1 || len(parts) > 3 {
		return false
	}

	// Check the hierarchy: city, province, country
	var city, province, country string
	if len(parts) == 3 {
		city = parts[0]
		province = parts[1]
		country = parts[2]
	} else if len(parts) == 2 {
		province = parts[0]
		country = parts[1]
	} else if len(parts) == 1 {
		country = parts[0]
	}

	// Validate the hierarchy
	if countryData, ok := Map[country]; ok {
		if province == "" {
			// Country is valid, no need to check province or city
			return true
		}
		if provinceData, exists := countryData[province]; exists {
			if city == "" {
				// Province is valid, no need to check city
				return true
			}
			// Check if the city exists in the province
			if _, exists := provinceData[city]; exists {
				return true
			}
		}
	}

	return false
}

// PrintMenu displays the menu options.
func PrintMenu() {
	fmt.Println("Distributor CLI\nCommands:\n1) Add distributor\n2) Add permission\n3) Link Distributors\n4) Remove distributor\n5) List distributor\n6) Remove Include/Exclude Permissions\n7) Get Parents\n8) Get Permissions\n9) Can Distribute\n10) Unlink Parent from Child\n- exit")
}

// UpperCaseAndTrimSpace converts a string to uppercase and removes extra spaces.
func UpperCaseAndTrimSpace(input string) string {
	return strings.ToUpper(strings.TrimSpace(input))
}

// AskUserForInput prompts the user for input with a message.
func AskUserForInput(userMessage string, reader *bufio.Reader) string {
	fmt.Println(userMessage)
	choice, _ := reader.ReadString('\n')
	return UpperCaseAndTrimSpace(choice)
}

// MakeLocationObject creates a LocationObject from a string input.
func MakeLocationObject(location string) LocationObject {
	array := strings.Split(location, ",")
	obj := LocationObject{}

	arrayLen := len(array)

	switch arrayLen {
	case 3:
		obj.City = UpperCaseAndTrimSpace(array[0])
		obj.Province = UpperCaseAndTrimSpace(array[1])
		obj.Country = UpperCaseAndTrimSpace(array[2])
	case 2:
		obj.Province = UpperCaseAndTrimSpace(array[0])
		obj.Country = UpperCaseAndTrimSpace(array[1])
	case 1:
		obj.Country = UpperCaseAndTrimSpace(array[0])
	default:
		fmt.Println("Invalid location details")
	}

	return obj
}

// DumpMapToFile dumps the map to a specified file in JSON format
func DumpMapToFile() error {
	// Create the file
	file, err := os.Create("output.json")
	if err != nil {
		return err
	}
	defer file.Close()

	// Create a JSON encoder
	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ") // Optional: for pretty printing

	// Encode the map into the file
	if err := encoder.Encode(Map); err != nil {
		return err
	}

	return nil
}
