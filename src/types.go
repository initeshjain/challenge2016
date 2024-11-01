package challenge2016

type Distributor struct {
	Name        string
	Parent      *Distributor
	Permissions Permission // Permissions include both Include and Exclude
}

type Permission struct {
	Include map[string]map[string]map[string]bool // Include map with nested structure
	Exclude map[string]map[string]map[string]bool // Exclude map with nested structure
}

var Distributors map[string]*Distributor

var Map map[string]map[string]map[string]bool

type LocationObject struct {
	Country string
	Province string
	City string
}