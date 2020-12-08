package maidenhead

import (
	"testing"
)

var tests = []struct {
	point Point
	loc   string
	loc4  string
}{
	{Point{48.14666, 11.60833}, "JN58TD", "JN58TD25"},
	{Point{-34.91, -56.21166}, "GF15VC", "GF15VC41"},
	{Point{38.92, -77.065}, "FM18LW", "FM18LW20"},
	{Point{-41.28333, 174.745}, "RE78IR", "RE78IR92"},
	{Point{41.714775, -72.727260}, "FN31PR", "FN31PR21"},
	{Point{37.413708, -122.1073236}, "CM87WJ", "CM87WJ79"},
	{Point{35.0542, -85.1142}, "EM75KB", "EM75KB63"},
}

func TestGridSquare(t *testing.T) {
	for _, test := range tests {
		enc, err := test.point.GridSquare()
		if err != nil {
			t.Fatal(err)
		}
		if enc != test.loc {
			t.Fatalf("%s want %q, got %q\n", test.point, test.loc, enc)
		}
		t.Logf("%s encoded to %q\n", test.point, enc)
	}
}

func TestExtendedSquarePrecision(t *testing.T) {
	for _, test := range tests {
		got, err := test.point.Locator(ExtendedSquarePrecision)
		if err != nil {
			t.Fatal(err)
		}
		if got != test.loc4 {
			t.Fatalf("%s want %q, got %q\n", test.point, test.loc4, got)
		}
		t.Logf("%s encoded to %q\n", test.point, got)
	}
}

// invalid Maiden Head locators must return error
func TestParseInvalidLocatorStrict(t *testing.T) {
	locs := []string{
		"JN58td",
		"JN58TDAA",
		"JNH",
		"QN58jh",
		"JN77ya",
		" ",
		"JN55J",
		"JN89HA11aa2",
		"JN89HA11aa22",
	}

	for _, l := range locs {
		_, err := ParseLocatorStrict(l)
		if err == nil {
			t.Errorf("Parsing invalid locator '%s' with ParseLocatorStrict() doesn't return any error", l)
		} else {
			t.Logf("Parsing invalid locator '%s' returns error: %s", l, err)
		}
	}
}

// Distance between corner point and center of the locator square
func TestParseLocatorCentered(t *testing.T) {
	tests := []struct {
		loc          string
		distExpected float64
	}{
		{"JN89", 91.42870273454076},
		{"JN89HF", 3.8111046375990782},
		{"JN89HF23", 0.38109528459829756},
		{"JN89HF23ag", 0.015878904160500258},
	}

	for _, test := range tests {
		p, _ := ParseLocator(test.loc)
		pc, _ := ParseLocatorCentered(test.loc)

		dist := pc.Distance(p)

		if dist != test.distExpected {
			t.Errorf("Distance between the center and corner of square locator '%s' is %g, expected %g",
				test.loc, dist, test.distExpected)
		}
	}
}
