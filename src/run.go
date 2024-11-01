package challenge2016

import (
	"bufio"
	"fmt"
	"os"
)

func RunInteractive() {
	reader := bufio.NewReader(os.Stdin)

	// Init Distributors
	Distributors = make(map[string]*Distributor)

	PrintMenu()

	for {
		fmt.Print("> ")
		input, _ := reader.ReadString('\n')
		input = UpperCaseAndTrimSpace(input)

		if input == "EXIT" {
			fmt.Println("Exiting...")
			break
		}

		switch input {
		case "1":
			AddDistributor(reader)
		case "2":
			AddPermission(reader)
		case "3":
			LinkDistributor(reader)
		case "4":
			RemoveDistributor(reader)
		case "5":
			ListDistributors()
		case "6":
			WipePermissions(reader)
		case "7":
			GetParentChain(reader)
		case "8":
			GetPermissions(reader)
		case "9":
			CanDistribute(reader)
		case "10":
			UnlinkParent(reader)
		default:
			fmt.Println("Unknown command.")
			PrintMenu()
		}
	}
}
