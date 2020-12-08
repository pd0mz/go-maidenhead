// Package maidenhead implements the Maidenhead Locator System, a geographic
// coordinate system used by amataur radio (HAM) operators.
package maidenhead

import (
	"errors"
	"fmt"
	"math"
	"strings"
)

// Precision of the computed locator.
const (
	FieldPrecision = iota + 1
	SquarePrecision
	SubSquarePrecision
	ExtendedSquarePrecision
)

var (
	upper = "ABCDEFGHIJKLMNOPQRSTUVWX"
	lower = "abcdefghijklmnopqrstuvwx"
	digit = "0123456789"
)

// locator computes the Maidenhead Locator for a given position.
func locator(p Point, precision int) (string, error) {
	if math.IsNaN(p.Latitude) {
		return "", errors.New("maidenhead: latitude is not a digit")
	}
	if math.IsInf(p.Latitude, 0) {
		return "", errors.New("maidenhead: latitude is infinite")
	}
	if math.IsNaN(p.Longitude) {
		return "", errors.New("maidenhead: longitude is not a digit")
	}
	if math.IsInf(p.Longitude, 0) {
		return "", errors.New("maidenhead: longitude is infinite")
	}
	if math.Abs(p.Latitude) == 90.0 {
		return "", errors.New("maidenhead: grid square invalid at poles")
	} else if math.Abs(p.Latitude) > 90.0 {
		return "", fmt.Errorf("maidenhead: invalid latitude %.04f", p.Latitude)
	} else if math.Abs(p.Longitude) > 180.0 {
		return "", fmt.Errorf("maidenhead: invalid longitude %.05f", p.Longitude)
	}

	var (
		lat = p.Latitude + 90.0
		lng = p.Longitude + 180.0
		loc string
	)

	lat = lat/10.0 + 0.0000001
	lng = lng/20.0 + 0.0000001
	loc = loc + string(upper[int(lng)]) + string(upper[int(lat)])
	if precision == 1 {
		return loc, nil
	}
	lat = 10 * (lat - math.Floor(lat))
	lng = 10 * (lng - math.Floor(lng))
	loc = loc + fmt.Sprintf("%d%d", int(lng)%10, int(lat)%10)
	if precision == 2 {
		return loc, nil
	}
	lat = 24 * (lat - math.Floor(lat))
	lng = 24 * (lng - math.Floor(lng))
	loc = loc + string(upper[int(lng)]) + string(upper[int(lat)])
	if precision == 3 {
		return loc, nil
	}
	lat = 10 * (lat - math.Floor(lat))
	lng = 10 * (lng - math.Floor(lng))
	loc = loc + fmt.Sprintf("%d%d", int(lng)%10, int(lat)%10)
	if precision == 4 {
		return loc, nil
	}
	lat = 24 * (lat - math.Floor(lat))
	lng = 24 * (lng - math.Floor(lng))
	loc = loc + string(lower[int(lng)]) + string(lower[int(lat)])
	return loc, nil
}

var parseLocatorMult = []struct {
	s, p string
	mult float64
}{
	{upper[:18], lower[:18], 20.0},
	{upper[:18], lower[:18], 10.0},
	{digit[:10], digit[:10], 20.0 / 10.0},
	{digit[:10], digit[:10], 10.0 / 10.0},
	{upper[:24], lower[:24], 20.0 / (10.0 * 24.0)},
	{upper[:24], lower[:24], 10.0 / (10.0 * 24.0)},
	{digit[:10], digit[:10], 20.0 / (10.0 * 24.0 * 10.0)},
	{digit[:10], digit[:10], 10.0 / (10.0 * 24.0 * 10.0)},
	{lower[:24], lower[:24], 20.0 / (10.0 * 24.0 * 10.0 * 24.0)},
	{lower[:24], lower[:24], 10.0 / (10.0 * 24.0 * 10.0 * 24.0)},
}

var maxLocatorLength = len(parseLocatorMult)

func parseLocator(locator string, strict bool, centered bool) (point Point, err error) {
	var (
		lnglat = [2]float64{
			-180.0,
			-90.0,
		}
		i, j int
		char rune
	)

	if len(locator) > maxLocatorLength {
		err = fmt.Errorf("maidenhead: locator is too long (%d characters, maximum %d characters allowed)",
			len(locator), maxLocatorLength)
		return
	}

	if len(locator)%2 != 0 {
		err = fmt.Errorf("maidenhead: locator has odd number of characters")
		return
	}

	if strict {
		for i, char = range locator {
			if j = strings.Index(parseLocatorMult[i].s, string(char)); j < 0 {
				err = fmt.Errorf("maidenhead: invalid character at offset %d", i)
				return
			}
			lnglat[i%2] += float64(j) * parseLocatorMult[i].mult
		}
	} else {
		for i, char = range strings.ToLower(locator) {
			if j = strings.Index(parseLocatorMult[i].p, string(char)); j < 0 {
				err = fmt.Errorf("maidenhead: invalid character at offset %d", i)
				return
			}
			lnglat[i%2] += float64(j) * parseLocatorMult[i].mult
		}
	}

	if centered {
		lnglat[0] += parseLocatorMult[i-1].mult / 2.0
		lnglat[1] += parseLocatorMult[i].mult / 2.0
	}

	point = NewPoint(lnglat[1], lnglat[0])
	return
}
