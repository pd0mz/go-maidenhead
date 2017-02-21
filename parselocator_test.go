package maidenhead

import "testing"

// parsed locator must be translated to the same locator
// using GridSquare()
func TestParseLocator(t *testing.T) {
	var locTests = []string{
		"JN88RT",
		"JN89HF",
		"JN58TD",
		"GF15VC",
		"FM18LW",
		"RE78IR",
	}

	for _, loc := range locTests {

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

		if l != loc {
			t.Errorf("%s: parsed to %v produces %s\n", loc, point, l)
		}
	}
}
