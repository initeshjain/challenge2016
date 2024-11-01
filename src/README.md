# Challenge 2016 Distributor Management System

## Overview

The Challenge 2016 Distributor Management System is a CLI application designed to manage distributors, their hierarchical relationships, and permissions based on geographical locations. It provides functionality to load regions from a CSV file, link distributors, and check their permissions for specific locations.

## Data Structures

### Location Map

The application uses a hierarchical map structure to store geographical locations:

```go
map[string]map[string]map[string]bool
```

- **Country**: The first-level key representing a country.
- **Province**: The second-level key representing a province within the country.
- **City**: The third-level key representing a city within the province, with a boolean value indicating its inclusion.

### Distributor

The `Distributor` structure contains the following fields:

- **Name**: A string representing the distributor's name.
- **Permissions**: An object containing:
  - **Include**: A map where keys are countries, and values are maps of provinces (keys) to maps of cities (keys) with boolean values indicating inclusion.
  - **Exclude**: Similar to **Include**, but for excluded permissions.
- **Parent**: A pointer to another `Distributor`, representing the parent in the hierarchy.

### LocationObject

The `LocationObject` structure is used to represent a location with the following fields:

- **City**: The city name (string).
- **Province**: The province name (string).
- **Country**: The country name (string).

## Functions

### LoadRegions

```go
func LoadRegions(filename string) error
```

- **Purpose**: Loads regions from a specified CSV file and populates the location map.
- **Parameters**: 
  - `filename`: The path to the CSV file containing region data.
- **Returns**: An error if loading fails, or nil on success.
- **Implementation Details**: 
  - Opens the CSV file and initializes the location map.
  - Reads each row, extracting country, province, and city names.
  - Updates the map structure to reflect the hierarchy of locations.

### IsValidLocation

```go
func IsValidLocation(location string) bool
```

- **Purpose**: Validates a location string against the populated location map.
- **Parameters**: 
  - `location`: A string representing a location, which can be in the format of "City, Province, Country".
- **Returns**: `true` if the location is valid; otherwise, `false`.
- **Implementation Details**: 
  - Splits the location into components and checks if each part exists in the location map according to the hierarchy (country > province > city).

### PrintMenu

```go
func PrintMenu()
```

- **Purpose**: Displays the CLI menu options for user interaction.
- **Implementation Details**: 
  - Prints the available commands to the console.

### UpperCaseAndTrimSpace

```go
func UpperCaseAndTrimSpace(input string) string
```

- **Purpose**: Converts a string to uppercase and trims leading/trailing whitespace.
- **Parameters**: 
  - `input`: The string to be processed.
- **Returns**: The processed string.
- **Implementation Details**: 
  - Utilizes standard string functions to format the input.

### AskUserForInput

```go
func AskUserForInput(userMessage string, reader *bufio.Reader) string
```

- **Purpose**: Prompts the user for input and reads their response.
- **Parameters**: 
  - `userMessage`: The prompt message to display to the user.
  - `reader`: A buffered reader for reading user input.
- **Returns**: The user's input after processing.
- **Implementation Details**: 
  - Prints the message and reads input from the console, applying the `UpperCaseAndTrimSpace` function.

### MakeLocationObject

```go
func MakeLocationObject(location string) LocationObject
```

- **Purpose**: Creates a `LocationObject` from a formatted location string.
- **Parameters**: 
  - `location`: A string representing a location in the format "City, Province, Country".
- **Returns**: A `LocationObject` populated with the extracted city, province, and country.
- **Implementation Details**: 
  - Splits the input string by commas and assigns values to the `LocationObject` fields based on the number of components.

### LinkDistributor

```go
func LinkDistributor(reader *bufio.Reader)
```

- **Purpose**: Links a child distributor to a parent distributor.
- **Parameters**: 
  - `reader`: A buffered reader for reading user input.
- **Implementation Details**: 
  - Prompts for parent and child distributor names.
  - Validates the existence of both distributors.
  - Ensures the child distributor does not have existing permissions.
  - Checks for potential cycles in the distributor hierarchy to prevent linking loops.

### AddDistributor

```go
func AddDistributor(reader *bufio.Reader)
```

- **Purpose**: Adds a new distributor to the `Distributors` map.
- **Parameters**: 
  - `reader`: A buffered reader for capturing user input.
- **Returns**: None.
- **Implementation Details**: 
  - Prompts the user for the distributor's name and checks if it already exists in the `Distributors` map. If it exists, a message is displayed; otherwise, a new distributor is created and added to the map.

### AddPermission

```go
func AddPermission(reader *bufio.Reader)
```

- **Purpose**: Adds location permissions (include/exclude) to a specified distributor.
- **Parameters**: 
  - `reader`: A buffered reader for capturing user input.
- **Returns**: None.
- **Implementation Details**: 
  - Retrieves the distributor by name, prompts for permission type (include or exclude), and requests location details. Validates the location and checks existing permissions before adding the new permission. If moving a location from include to exclude or vice versa, user consent is sought.

### GetPermissions

```go
func GetPermissions(reader *bufio.Reader)
```

- **Purpose**: Displays the permissions for a specified distributor.
- **Parameters**: 
  - `reader`: A buffered reader for capturing user input.
- **Returns**: None.
- **Implementation Details**: 
  - Prompts for the distributor's name, retrieves the distributor from the map, and iterates through its included and excluded permissions, printing them to the console.

### GetParentChain

```go
func GetParentChain(reader *bufio.Reader)
```

- **Purpose**: Displays the hierarchy of parent distributors for a specified distributor.
- **Parameters**: 
  - `reader`: A buffered reader for capturing user input.
- **Returns**: None.
- **Implementation Details**: 
  - Prompts for the distributor's name and retrieves the distributor. It traverses up the hierarchy of parent distributors, collecting their names, and prints them in order from the topmost parent to the specified distributor.

### ListDistributors

```go
func ListDistributors()
```

- **Purpose**: Lists all distributors in the system.
- **Parameters**: None.
- **Returns**: None.
- **Implementation Details**: 
  - Iterates through the `Distributors` map and prints the names of all distributors.

### RemoveDistributor

```go
func RemoveDistributor(reader *bufio.Reader)
```

- **Purpose**: Removes a distributor from the `Distributors` map if it has no children.
- **Parameters**: 
  - `reader`: A buffered reader for capturing user input.
- **Returns**: None.
- **Implementation Details**: 
  - Prompts for the distributor's name and checks if it exists. If it has child distributors (i.e., other distributors linked to it as parents), it cannot be removed, and an error message is shown. If no children exist, the distributor is deleted from the map.

### WipePermissions

```go
func WipePermissions(reader *bufio.Reader)
```

- **Purpose**: Clears all include or exclude permissions for a specified distributor.
- **Parameters**: 
  - `reader`: A buffered reader for capturing user input.
- **Returns**: None.
- **Implementation Details**: 
  - Prompts for the type of permissions to wipe (include or exclude) and the distributor's name. It checks if the distributor exists and then resets the relevant permissions map to a new empty map, effectively clearing the permissions.

### CanDistribute

```go
func CanDistribute(reader *bufio.Reader)
```

- **Purpose**: Checks if a specified distributor can distribute in a given location.
- **Parameters**: 
  - `reader`: A buffered reader for capturing user input.
- **Returns**: None.
- **Implementation Details**: 
  - Prompts for the distributor's name and location details. It creates a

 `LocationObject` and checks if the distributor has permission to operate in the specified location.
