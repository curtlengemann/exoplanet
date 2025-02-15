package main

import (
	"fmt"
	"main/internal/exoplanetCatalog"
)

// main kicks off the program
func main() {
	err := exoplanetCatalog.DisplayExoplanetData()
	if err != nil {
		fmt.Println("Error fetching exoplanets:", err.Error())
	}
}
