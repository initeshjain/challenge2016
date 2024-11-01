package challenge2016

import (
	"bufio"
	"fmt"
)

// AddDistributor adds a new distributor to the distributors map.
func AddDistributor(reader *bufio.Reader) {
	name := AskUserForInput("Enter name of distributor > ", reader)

	if _, exists := Distributors[name]; exists {
		fmt.Printf("Distributor %s already exists.\n", name)
		return
	}

	Distributors[name] = NewDistributor(name, nil)
	fmt.Printf("Added new distributor: %s\n", name)
}

func AddPermission(reader *bufio.Reader) {
	name := AskUserForInput("Enter name of distributor > ", reader)
	distributor, exists := Distributors[name]
	if !exists {
		fmt.Println("Distributor not found.")
		return
	}

	choice := AskUserForInput("Select permission type\n1) Include\n2) Exclude", reader)
	permissions := AskUserForInput("Write location details in the format <City>,<Province>,<Country>", reader)

	if !IsValidLocation(permissions) {
		fmt.Printf("Location %s is not valid.\n", permissions)
		return
	}

	locationObj := MakeLocationObject(permissions)

	switch choice {
	case "1": // Include
		if distributor.Excludes(locationObj) {
			msg := fmt.Sprintf("Location %s already excluded for distributor %s. Do you want to move it from Excluded to Included? Type yes or no > ", permissions, distributor.Name)
			consent := AskUserForInput(msg, reader)
			if consent != "YES" {
				return
			}
		}
		distributor.RemoveLocation(locationObj, "EXCLUDE")

		if distributor.Parent != nil && !distributor.Parent.Includes(locationObj) {
			fmt.Println("Cannot add as parent has not allowed it to include/exclude yet.")
		} else {
			distributor.AddInclude(locationObj)
			fmt.Printf("Added %s to the include permissions for distributor %s.\n", permissions, distributor.Name)
		}

	case "2": // Exclude
		if distributor.Includes(locationObj) {
			msg := fmt.Sprintf("Location %s already included for distributor %s. Do you want to move it from Included to Excluded? Type yes or no > ", permissions, distributor.Name)
			consent := AskUserForInput(msg, reader)
			if consent != "YES" {
				return
			}
		}
		distributor.RemoveLocation(locationObj, "INCLUDE")

		if distributor.Parent != nil && !distributor.Parent.Includes(locationObj) {
			fmt.Println("Cannot add as parent has not allowed it to include/exclude yet.")
		} else {
			distributor.AddExclude(locationObj)
			fmt.Printf("Added %s to the exclude permissions for distributor %s.\n", permissions, distributor.Name)
		}
	default:
		fmt.Println("Invalid choice. Please select 1 or 2.")
	}
}

// LinkDistributor links a child distributor to a parent distributor.
func LinkDistributor(reader *bufio.Reader) {
	parentName := AskUserForInput("Write Parent Distributor Name > ", reader)
	childName := AskUserForInput("Write Child Distributor Name > ", reader)

	parent, parentExists := Distributors[parentName]
	if !parentExists {
		fmt.Printf("Parent distributor %s not found.\n", parentName)
		return
	}

	child, childExists := Distributors[childName]
	if !childExists {
		fmt.Printf("Child distributor %s not found.\n", childName)
		return
	}

	if len(child.Permissions.Include) > 0 || len(child.Permissions.Exclude) > 0 {
		fmt.Printf("Distributor %s has existing permissions. Please wipe all permissions and try again.\n", childName)
		return
	}

	// Check if Child is already a parent to current parent (Preventing Cycle)
	if parent.Parent == child {
		fmt.Printf("%s is already a parent of %s. This operation would create a cycle. Please unlink it first.\n", childName, parentName)
		return
	}

	child.Parent = parent
	fmt.Printf("Linked distributor %s to parent %s\n", childName, parentName)
}

// GetPermissions displays the permissions for a distributor
func GetPermissions(reader *bufio.Reader) {
	name := AskUserForInput("Write distributor name > ", reader)

	if distributor, exists := Distributors[name]; exists {
		fmt.Printf("Permissions for %s:\n", name)
		for countryName, provinceMap := range distributor.Permissions.Include {
			for provinceName, cityMap := range provinceMap {
				for cityName := range cityMap {
					fmt.Printf("INCLUDE: %s, %s, %s\n", cityName, provinceName, countryName)
				}
			}
		}

		for countryName, provinceMap := range distributor.Permissions.Exclude {
			for provinceName, cityMap := range provinceMap {
				for cityName := range cityMap {
					fmt.Printf("EXCLUDE: %s, %s, %s\n", cityName, provinceName, countryName)
				}
			}
		}
	} else {
		fmt.Println("Distributor not found.")
	}
}

// GetParentChain displays the parent chain for a distributor
func GetParentChain(reader *bufio.Reader) {
	name := AskUserForInput("Write distributor name > ", reader)
	if distributor, exists := Distributors[name]; exists {
		var parentChain []string
		parent := distributor.Parent

		for parent != nil {
			parentChain = append(parentChain, parent.Name)
			parent = parent.Parent
		}

		for i := len(parentChain) - 1; i >= 0; i-- {
			if i == len(parentChain)-1 {
				fmt.Print(parentChain[i])
			} else {
				fmt.Print(" < ", parentChain[i])
			}
		}
		fmt.Println()
	} else {
		fmt.Println("Distributor not found.")
	}
}

// ListDistributors lists all distributors
func ListDistributors() {
	fmt.Println("Distributors: ")
	for name := range Distributors {
		fmt.Println("-", name)
	}
}

// RemoveDistributor removes a distributor if it has no children
func RemoveDistributor(reader *bufio.Reader) {
	distributorName := AskUserForInput("Write Distributor name to delete", reader)

	if _, exists := Distributors[distributorName]; !exists {
		fmt.Printf("Distributor %s not found.\n", distributorName)
		return
	}

	for name, dist := range Distributors {
		if dist.Parent != nil && dist.Parent.Name == distributorName {
			fmt.Printf("Distributor %s is a parent of %s. Remove all children from %s and try again.\n", distributorName, name, distributorName)
			return
		}
	}

	delete(Distributors, distributorName)
	fmt.Printf("Removed %s successfully from distributors list.\n", distributorName)
}

// WipePermissions clears the include/exclude list for a distributor
func WipePermissions(reader *bufio.Reader) {
	subcmd := AskUserForInput("Which permission you want to delete\n1) Include\n2) Exclude", reader)

	distributorName := AskUserForInput("Write name of distributor to delete its permissions", reader)

	distributor, exists := Distributors[distributorName]
	if !exists {
		fmt.Printf("Distributor %s not found.\n", distributorName)
		return
	}

	switch subcmd {
	case "1": // Include
		distributor.Permissions.Include = make(map[string]map[string]map[string]bool)
		fmt.Printf("Wiped INCLUDE list for distributor %s.\n", distributorName)
	case "2": // Exclude
		distributor.Permissions.Exclude = make(map[string]map[string]map[string]bool)
		fmt.Printf("Wiped EXCLUDE list for distributor %s.\n", distributorName)
	default:
		fmt.Println("Unknown command. Please select 1 for Include or 2 for Exclude.")
	}
}

// Check if distributor can distribute in given location
func CanDistribute(reader *bufio.Reader) {
	name := AskUserForInput("Enter name of distributor > ", reader)
	distributor, exists := Distributors[name]
	if !exists {
		fmt.Println("Distributor not found.")
		return
	}

	permission := AskUserForInput("Write location detail in the format <City>,<Province>,<Country>", reader)

	permission = UpperCaseAndTrimSpace(permission)
	locationObject := MakeLocationObject(permission)
	result := distributor.CanDistribute(locationObject)

	if result {
		fmt.Printf("Distributor %s can distribute in %s\n", distributor.Name, permission)
	} else {
		fmt.Printf("Distributor %s cannot distribute in %s\n", distributor.Name, permission)
	}
}

// Unlink Parent
func UnlinkParent(reader *bufio.Reader) {
	name := AskUserForInput("Enter name of child distributor > ", reader)
	distributor, exists := Distributors[name]
	if !exists {
		fmt.Println("Distributor not found.")
		return
	}

	parent := distributor.Parent

	if distributor.Parent != nil {
		distributor.Parent = nil
		fmt.Printf("Distributor %s unlinked from its Parent %s\n", distributor.Name, parent.Name)
	} else {
		fmt.Printf("Distributor %s is Orphan\n", distributor.Name)
	}
}
