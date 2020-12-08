package maidenhead

import (
	"math"
	"testing"
)

var pointTests = []struct {
	point   Point
	bearing float64
	compass string
}{
	{Point{48.14666, 11.60833}, 195, "SSW"},
	{Point{-34.91, -56.21166}, 69, "ENE"},
	{Point{38.92, -77.065}, 98, "E"},
	{Point{-41.28333, 174.745}, 187, "S"},
	{Point{41.714775, -72.727260}, 101, "ESE"},
	{Point{37.413708, -122.1073236}, 69, "ENE"},
	{Point{35.0542, -85.1142}, 92, "E"},
}

func TestPointBearing(t *testing.T) {
	var center = NewPoint(0, 0)
	for _, test := range pointTests {
		bearing := math.Floor(test.point.Bearing(center))
		if bearing != test.bearing {
			t.Fatalf("%s -> %s, expected %0.f, got %0.f\n", test.point, center, test.bearing, bearing)
		}
		compass := test.point.CompassBearing(center)
		if compass != test.compass {
			t.Logf("%s -> %s, expected %q, got %q\n", test.point, center, test.compass, compass)
		}
		t.Logf("%s -> %s, bearing %.0f %s\n", test.point, center, bearing, compass)
	}
}
