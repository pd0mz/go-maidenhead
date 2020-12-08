package maidenhead

import (
	"fmt"
	"math"
)

const (
	// Earth radius
	r = 6371
)

var compassBearing = []struct {
	label        string
	start, ended float64
}{
	{"N", 000.00, 011.25}, {"NNE", 011.25, 033.75}, {"NE", 033.75, 056.25}, {"ENE", 056.25, 078.75},
	{"E", 078.75, 101.25}, {"ESE", 101.25, 123.75}, {"SE", 123.75, 146.25}, {"SSE", 146.25, 168.75},
	{"S", 168.75, 191.25}, {"SSW", 191.25, 213.75}, {"SW", 213.75, 236.25}, {"WSW", 236.25, 258.75},
	{"W", 258.75, 281.25}, {"WNW", 281.25, 303.75}, {"NW", 303.75, 326.25}, {"NNW", 326.25, 348.75},
	{"N", 348.75, 360.00},
}

// Point is a geographical point on the map.
type Point struct {
	Latitude  float64
	Longitude float64
}

// NewPoint returns a new Point structure with given latitude and longitude.
func NewPoint(latitude, longitude float64) Point {
	return Point{latitude, longitude}
}

// ParseLocator parses a Maidenhead Locator with permissive rule matching.
func ParseLocator(locator string) (Point, error) {
	return parseLocator(locator, false, false)
}

// ParseLocatorStrict parses a Maidenhead Locator with strict rule matching.
func ParseLocatorStrict(locator string) (Point, error) {
	return parseLocator(locator, true, false)
}

// ParseLocatorCentered parses a Maidenhead Locator with permissive rule matching.
// Returns Points structure with coordinates of the square center
func ParseLocatorCentered(locator string) (Point, error) {
	return parseLocator(locator, false, true)
}

// ParseLocatorStrictCentered parses a Maidenhead Locator with strict rule matching.
// Returns Points structure with coordinates of the square center
func ParseLocatorStrictCentered(locator string) (Point, error) {
	return parseLocator(locator, true, true)
}

// EqualTo returns true if the coordinates point to the same geographical location.
func (p Point) EqualTo(other Point) bool {
	var (
		dlat = p.Latitude - other.Latitude
		dlng = p.Longitude - other.Longitude
	)

	for dlat < -180.0 {
		dlat += 360.0
	}
	for dlat > 180.0 {
		dlat -= 360.0
	}
	for dlng < -90.0 {
		dlng += 90.0
	}
	for dlng > 90.0 {
		dlng -= 90.0
	}

	return dlat == 0.0 && dlng == 0.0
}

// Bearing calculates the (approximate) bearing to another heading.
func (p Point) Bearing(heading Point) float64 {
	var (
		hn = p.Latitude / 180 * math.Pi
		he = p.Longitude / 180 * math.Pi
		n  = heading.Latitude / 180 * math.Pi
		e  = heading.Longitude / 180 * math.Pi
		co = math.Cos(he-e)*math.Cos(hn)*math.Cos(n) + math.Sin(hn)*math.Sin(n)
		ca = math.Atan(math.Abs(math.Sqrt(1-co*co) / co))
	)

	if co < 0.0 {
		ca = math.Pi - ca
	}

	var si = math.Sin(e-he) * math.Cos(n) * math.Cos(hn)
	co = math.Sin(n) - math.Sin(hn)*math.Cos(ca)
	var az = math.Atan(math.Abs(si / co))
	if co < 0.0 {
		az = math.Pi - az
	}
	if si < 0.0 {
		az = -az
	}
	if az < 0.0 {
		az = az + 2.0*math.Pi
	}
	return az * 180 / math.Pi
}

// CompassBearing returns the compass bearing to a heading.
func (p Point) CompassBearing(heading Point) string {
	bearing := p.Bearing(heading)
	for bearing < 0.0 {
		bearing += 360.0
	}
	for bearing > 360.0 {
		bearing -= 360.0
	}

	for _, compass := range compassBearing {
		if bearing >= compass.start && bearing <= compass.ended {
			return compass.label
		}
	}

	// Should never reach
	return ""
}

// Distance calculates the (approximate) distance to another point in km.
func (p Point) Distance(other Point) float64 {
	var (
		hn = p.Latitude / 180 * math.Pi
		he = p.Longitude / 180 * math.Pi
		n  = other.Latitude / 180 * math.Pi
		e  = other.Longitude / 180 * math.Pi
		co = math.Cos(he-e)*math.Cos(hn)*math.Cos(n) + math.Sin(hn)*math.Sin(n)
		ca = math.Atan(math.Abs(math.Sqrt(1-co*co) / co))
	)

	if co < 0.0 {
		ca = math.Pi - ca
	}

	return r * ca
}

// GridSquare returns a Maidenhead Locator for the point coordinates.
func (p Point) GridSquare() (string, error) {
	return locator(p, SubSquarePrecision)
}

// Locator returns a Maidenhead Locator for the point coordinates with
// specified precision
func (p Point) Locator(precision int) (string, error) {
	return locator(p, precision)
}

// String returns a stringified Point structure.
func (p Point) String() string {
	return fmt.Sprintf("Point(%f, %f)", p.Latitude, p.Longitude)
}
