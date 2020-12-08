package maidenhead

import (
	"math"
	"strings"
	"testing"
)

// parsed locator must be translated to the same locator
// using GridSquare()
func TestParseLocator(t *testing.T) {
	var locTests = map[string]Point{
		"JN88RT": Point{48.791667, 17.416667},
		"JN89HF": Point{49.208333, 16.583333},
		"JN58TD": Point{48.125000, 11.583333},
		"GF15VC": Point{-34.916667, -56.250000},
		"FM18LW": Point{38.916667, -77.083333},
		"RE78IR": Point{-41.291667, 174.666667},
		"PM45lm": Point{35.5, 128.916667},
	}

	for loc, p := range locTests {
		point, err := ParseLocator(loc)
		if err != nil {
			t.Errorf("%s parsing error: %s", loc, err)
			continue
		}

		l, err := point.GridSquare()
		if err != nil {
			t.Errorf("%s: %v to GridSquare(): %s", loc, point, err)
			continue
		}

		if !strings.EqualFold(l, loc) {
			t.Errorf("%s: parsed to %v produces %s\n", loc, point, l)
		}

		if !(almostEqual(p.Latitude, point.Latitude) && almostEqual(p.Longitude, point.Longitude)) {
			t.Errorf("%s: at %s, expeted %s", loc, point, p)
		}
	}
}

func almostEqual(a, b float64) bool {
	return math.Abs(a-b) < 1e-06
}
