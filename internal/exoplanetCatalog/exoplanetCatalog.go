package exoplanetCatalog

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"slices"
	"strings"
)

// typeFlag is a custom type used to represent different classifications of planetary systems.
type typeFlag int

const (
	noKnownStellarBinaryCompanion typeFlag = iota
	pTypeBinary
	sTypeBinary
	orphanPlanet
)

// exoplanet defines a struct based on the given JSON structure
type exoplanet struct {
	PlanetIdentifier string        `json:"PlanetIdentifier"`
	TypeFlag         typeFlag      `json:"TypeFlag"`
	RadiusJpt        nillableFloat `json:"RadiusJpt"`
	DiscoveryYear    nillableFloat `json:"DiscoveryYear"`
	HostStarTempK    nillableFloat `json:"HostStarTempK"`
}

// URL for the exoplanet data
const dataSourceURL = "https://gist.githubusercontent.com/joelbirchler/66cf8045fcbb6515557347c05d789b4a/raw/9a196385b44d4288431eef74896c0512bad3defe/exoplanets"

func DisplayExoplanetData() error {
	exoplanets, err := fetchExoplanets(dataSourceURL)
	if err != nil {
		return err
	}

	fmt.Println("Number of orphan planets:", getNumberOfOrphanedPlanets(exoplanets))
	fmt.Println("Planet(s) orbiting hottest star(s):", strings.Join(getPlanetsOrbitingHottestStar(exoplanets), ", "))
	printTimeline(getTimeline(exoplanets))
	return nil
}

// fetchExoplanets retrieves a list of exoplanets from a remote URL.
// It sends an HTTP GET request to a predefined URL, reads the response body,
// and unmarshals the JSON data into a slice of Exoplanet models.
func fetchExoplanets(url string) ([]exoplanet, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var exoplanets []exoplanet
	err = json.Unmarshal(body, &exoplanets)
	if err != nil {
		return nil, err
	}
	return exoplanets, nil
}

// getNumberOfOrphanedPlanets returns the number of orphan planets using the exoplanet's type flag.
func getNumberOfOrphanedPlanets(exoplanets []exoplanet) int {
	numOrphanPlanets := 0
	for _, planet := range exoplanets {
		if planet.TypeFlag == orphanPlanet {
			numOrphanPlanets++
		}
	}
	return numOrphanPlanets
}

// getPlanetsOrbitingHottestStar returns the name(s) of the planet oribiting the hottest star(s).
func getPlanetsOrbitingHottestStar(exoplanets []exoplanet) []string {
	var hottestPlanets []string
	hottestPlanetTemp := float64(-1)

	for _, planet := range exoplanets {
		temp := planet.HostStarTempK.Value
		if temp == nil {
			continue
		}
		if hottestPlanetTemp < *temp {
			hottestPlanetTemp = *temp
			hottestPlanets = []string{planet.PlanetIdentifier}
		} else if hottestPlanetTemp == *temp {
			hottestPlanets = append(hottestPlanets, planet.PlanetIdentifier)
		}
	}
	return hottestPlanets
}

type planetGrouping struct {
	smallSizePlanets  int
	mediumSizePlanets int
	largeSizePlanets  int
}

// getTimeline returns the number of planets discovered per year grouped by size.
func getTimeline(exoplanets []exoplanet) map[int]planetGrouping {
	const smallSize, mediumSize = float64(1), float64(2)

	planetTimeline := make(map[int]planetGrouping, len(exoplanets))
	for _, planet := range exoplanets {
		year := planet.DiscoveryYear.Value
		if year == nil {
			continue
		}

		radius := planet.RadiusJpt.Value
		if radius == nil {
			continue
		}

		// This cast is safe as year is always an int
		intYear := int(*year)
		planetYear := planetTimeline[intYear]

		if *radius < smallSize {
			planetYear.smallSizePlanets++
		} else if *radius < mediumSize {
			planetYear.mediumSizePlanets++
		} else {
			planetYear.largeSizePlanets++
		}
		planetTimeline[intYear] = planetYear
	}
	return planetTimeline
}

// printTimeline takes an unordered map and sorts it to make a pretty printed timeline.
func printTimeline(timeline map[int]planetGrouping) {
	years := make([]int, 0, len(timeline))
	for year := range timeline {
		years = append(years, year)
	}
	slices.Sort(years)

	for _, year := range years {
		group := timeline[year]
		fmt.Println("In", year, "we discovered", group.smallSizePlanets, "small planets,", group.mediumSizePlanets, "medium planets, and", group.largeSizePlanets, "large planets.")
	}
}
