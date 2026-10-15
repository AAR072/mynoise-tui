package player

import (
	"encoding/json"
	"math"
	"strings"
	"testing"

	"github.com/aar072/mynoise-tui/classes"
)

func TestSoundJSEscapesNamesAndPreservesPrecision(t *testing.T) {
	sound := classes.Sound{Name: "quotes\" backslash\\ newline\n", Sliders: [10]float64{0.123456789}}
	code, err := soundJS(sound)
	if err != nil {
		t.Fatal(err)
	}
	encoded, _ := json.Marshal(sound.Name)
	if !strings.Contains(code, string(encoded)) || !strings.Contains(code, "0.123456789") {
		t.Fatalf("bad JS: %s", code)
	}
	sound.Sliders[0] = math.NaN()
	if _, err := soundJS(sound); err == nil {
		t.Fatal("expected invalid slider error")
	}
	code, err = soundJS(classes.DefaultSound)
	if err != nil || code != "resetSliders()" {
		t.Fatalf("default: %s, %v", code, err)
	}
}
