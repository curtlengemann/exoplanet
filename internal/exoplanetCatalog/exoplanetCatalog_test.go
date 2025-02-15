package exoplanetCatalog

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDisplayExoplanetData(t *testing.T) {

	got := captureOutput(func() {
		DisplayExoplanetData()
	})

	// Just assert that we print something. We don't care what it is.
	// This dataset could update tomorrow which would cause this test to break
	// if we assert exactly what we are printing.
	assert.True(t, len(got) > 0)
}

func TestFetchExoplanets(t *testing.T) {
	tests := []struct {
		name         string
		httpStatus   int
		responseBody string
		wantErr      bool
		expected     []exoplanet
	}{
		{
			name:         "ValidResponse",
			httpStatus:   http.StatusOK,
			responseBody: `[{"PlanetIdentifier":"Planet1"},{"PlanetIdentifier":"Planet2"}]`,
			wantErr:      false,
			expected:     []exoplanet{{PlanetIdentifier: "Planet1"}, {PlanetIdentifier: "Planet2"}},
		},
		{
			name:         "EmptyResponse",
			httpStatus:   http.StatusOK,
			responseBody: `[]`,
			wantErr:      false,
			expected:     []exoplanet{},
		},
		{
			name:         "InvalidJSON",
			httpStatus:   http.StatusOK,
			responseBody: `{`,
			wantErr:      true,
			expected:     nil,
		},
		{
			name:         "HTTPError",
			httpStatus:   http.StatusInternalServerError,
			responseBody: ``,
			wantErr:      true,
			expected:     nil,
		},
		{
			name:         "NilResponse",
			httpStatus:   http.StatusOK,
			responseBody: `null`,
			wantErr:      false,
			expected:     nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.httpStatus)
				w.Write([]byte(tt.responseBody))
			}))
			defer server.Close()

			resp, err := fetchExoplanets(server.URL)
			assert.Equal(t, tt.wantErr, err != nil)
			assert.Equal(t, tt.expected, resp)
		})
	}
}

func TestGetNumberOfOrphanedPlanets(t *testing.T) {
	tests := []struct {
		name       string
		exoplanets []exoplanet
		expected   int
	}{
		{
			name:       "NoOrphanPlanets",
			exoplanets: []exoplanet{{TypeFlag: pTypeBinary}},
			expected:   0,
		},
		{
			name:       "OneOrphanPlanet",
			exoplanets: []exoplanet{{TypeFlag: orphanPlanet}},
			expected:   1,
		},
		{
			name: "MultipleOrphanPlanets",
			exoplanets: []exoplanet{
				{TypeFlag: orphanPlanet},
				{TypeFlag: orphanPlanet},
			},
			expected: 2,
		},
		{
			name:       "NilInput",
			exoplanets: nil,
			expected:   0,
		},
		{
			name: "MixedPlanets",
			exoplanets: []exoplanet{
				{TypeFlag: pTypeBinary},
				{TypeFlag: sTypeBinary},
				{TypeFlag: orphanPlanet},
				{TypeFlag: noKnownStellarBinaryCompanion},
			},
			expected: 1,
		},
		{
			name: "NoTypeFlag",
			exoplanets: []exoplanet{
				{TypeFlag: 999}, // Invalid TypeFlag
			},
			expected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, getNumberOfOrphanedPlanets(tt.exoplanets))
		})
	}
}

func TestGetPlanetOrbitingHottestStar(t *testing.T) {
	tests := []struct {
		name       string
		exoplanets []exoplanet
		expected   []string
	}{
		{
			name:       "NoPlanets",
			exoplanets: nil,
			expected:   nil,
		},
		{
			name: "OnePlanet",
			exoplanets: []exoplanet{
				{PlanetIdentifier: "Planet1", HostStarTempK: nillableFloat{Value: ptrFloat64(1)}},
			},
			expected: []string{"Planet1"},
		},
		{
			name: "MultiplePlanets",
			exoplanets: []exoplanet{
				{PlanetIdentifier: "Planet1", HostStarTempK: nillableFloat{Value: ptrFloat64(1)}},
				{PlanetIdentifier: "Planet2", HostStarTempK: nillableFloat{Value: ptrFloat64(2)}},
			},
			expected: []string{"Planet2"},
		},

		{
			name: "MultiplePlanetsWithSameTemp",
			exoplanets: []exoplanet{
				{PlanetIdentifier: "Planet1", HostStarTempK: nillableFloat{Value: ptrFloat64(2)}},
				{PlanetIdentifier: "Planet2", HostStarTempK: nillableFloat{Value: ptrFloat64(2)}},
				{PlanetIdentifier: "Planet3", HostStarTempK: nillableFloat{Value: ptrFloat64(1)}},
			},
			expected: []string{"Planet1", "Planet2"},
		},
		{
			name: "NoTemperatureData",
			exoplanets: []exoplanet{
				{PlanetIdentifier: "Planet1", HostStarTempK: nillableFloat{Value: nil}},
			},
			expected: nil,
		},
		{
			name: "MixedTemperatureData",
			exoplanets: []exoplanet{
				{PlanetIdentifier: "Planet1", HostStarTempK: nillableFloat{Value: ptrFloat64(2)}},
				{PlanetIdentifier: "Planet2", HostStarTempK: nillableFloat{Value: nil}},
				{PlanetIdentifier: "Planet3", HostStarTempK: nillableFloat{Value: ptrFloat64(1)}},
			},
			expected: []string{"Planet1"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, getPlanetsOrbitingHottestStar(tt.exoplanets))
		})
	}
}

func TestGetTimeline(t *testing.T) {
	tests := []struct {
		name       string
		exoplanets []exoplanet
		expected   map[int]planetGrouping
	}{
		{
			name:       "NoPlanets",
			exoplanets: nil,
			expected:   map[int]planetGrouping{},
		},
		{
			name: "OneSmallPlanet",
			exoplanets: []exoplanet{
				{DiscoveryYear: nillableFloat{Value: ptrFloat64(2000)}, RadiusJpt: nillableFloat{Value: ptrFloat64(0.5)}},
			},
			expected: map[int]planetGrouping{
				2000: {smallSizePlanets: 1},
			},
		},
		{
			name: "OneMediumPlanet",
			exoplanets: []exoplanet{
				{DiscoveryYear: nillableFloat{Value: ptrFloat64(2001)}, RadiusJpt: nillableFloat{Value: ptrFloat64(1.5)}},
			},
			expected: map[int]planetGrouping{
				2001: {mediumSizePlanets: 1},
			},
		},
		{
			name: "OneLargePlanet",
			exoplanets: []exoplanet{
				{DiscoveryYear: nillableFloat{Value: ptrFloat64(2002)}, RadiusJpt: nillableFloat{Value: ptrFloat64(2.5)}},
			},
			expected: map[int]planetGrouping{
				2002: {largeSizePlanets: 1},
			},
		},
		{
			name: "MixedPlanets",
			exoplanets: []exoplanet{
				{DiscoveryYear: nillableFloat{Value: ptrFloat64(2000)}, RadiusJpt: nillableFloat{Value: ptrFloat64(0.5)}},
				{DiscoveryYear: nillableFloat{Value: ptrFloat64(2001)}, RadiusJpt: nillableFloat{Value: ptrFloat64(1.5)}},
				{DiscoveryYear: nillableFloat{Value: ptrFloat64(2002)}, RadiusJpt: nillableFloat{Value: ptrFloat64(2.5)}},
				{DiscoveryYear: nillableFloat{Value: ptrFloat64(2000)}, RadiusJpt: nillableFloat{Value: ptrFloat64(.9)}},
				{DiscoveryYear: nillableFloat{Value: ptrFloat64(2000)}, RadiusJpt: nillableFloat{Value: ptrFloat64(1)}},
				{DiscoveryYear: nillableFloat{Value: ptrFloat64(2000)}, RadiusJpt: nillableFloat{Value: ptrFloat64(2)}},
			},
			expected: map[int]planetGrouping{
				2000: {
					smallSizePlanets:  2,
					mediumSizePlanets: 1,
					largeSizePlanets:  1,
				},
				2001: {mediumSizePlanets: 1},
				2002: {largeSizePlanets: 1},
			},
		},
		{
			name: "NilYear",
			exoplanets: []exoplanet{
				{DiscoveryYear: nillableFloat{Value: nil}, RadiusJpt: nillableFloat{Value: ptrFloat64(0.5)}},
			},
			expected: map[int]planetGrouping{},
		},
		{
			name: "NilRadius",
			exoplanets: []exoplanet{
				{DiscoveryYear: nillableFloat{Value: ptrFloat64(2000)}, RadiusJpt: nillableFloat{Value: nil}},
			},
			expected: map[int]planetGrouping{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, getTimeline(tt.exoplanets))
		})
	}
}

func TestPrintTimeline(t *testing.T) {
	tests := []struct {
		name     string
		timeline map[int]planetGrouping
		expected string
	}{
		{
			name:     "Single entry",
			timeline: map[int]planetGrouping{2025: {smallSizePlanets: 2, mediumSizePlanets: 1, largeSizePlanets: 3}},
			expected: "In 2025 we discovered 2 small planets, 1 medium planets, and 3 large planets.\n",
		},
		{
			name: "Multiple entries out of order",
			timeline: map[int]planetGrouping{
				2024: {smallSizePlanets: 1, mediumSizePlanets: 2, largeSizePlanets: 1},
				2022: {smallSizePlanets: 3, mediumSizePlanets: 4, largeSizePlanets: 2},
				2023: {smallSizePlanets: 5, mediumSizePlanets: 6, largeSizePlanets: 7},
			},
			expected: "In 2022 we discovered 3 small planets, 4 medium planets, and 2 large planets.\n" +
				"In 2023 we discovered 5 small planets, 6 medium planets, and 7 large planets.\n" +
				"In 2024 we discovered 1 small planets, 2 medium planets, and 1 large planets.\n",
		},
		{
			name:     "No entries",
			timeline: map[int]planetGrouping{},
			expected: "",
		},
		{
			name:     "Nil map",
			timeline: nil,
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := captureOutput(func() {
				printTimeline(tt.timeline)
			})

			assert.Equal(t, tt.expected, got)
		})
	}
}

// Helper function to create a float64 pointer
func ptrFloat64(f float64) *float64 {
	return &f
}

// Helper function captures output from the print function
func captureOutput(f func()) string {
	r, w, _ := os.Pipe()
	defer r.Close()
	defer w.Close()

	stdout := os.Stdout
	defer func() { os.Stdout = stdout }()
	os.Stdout = w

	f()
	w.Close()

	var buf bytes.Buffer
	io.Copy(&buf, r)
	return buf.String()
}
