package maidenhead

import "testing"

var loc_tests = []string{
	"JN88RT",
	"JN89HF",
	"JN58TD",
	"GF15VC",
	"FM18LW",
	"RE78IR",
}

// parsed locator must be translated to the same locator
// using GridSquare()
func TestParseLocator(t *testing.T) {
	for _, loc := range loc_tests {

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
