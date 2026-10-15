package player

import (
	"encoding/json"
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/aar072/mynoise-tui/browser"
	"github.com/aar072/mynoise-tui/classes"
)

func PlayPreset(preset classes.Preset) error {
	if err := browser.NavigateTo(preset.Data.URL); err != nil {
		return err
	}
	return browser.CallJSFunction("resetSliders()")
}

func PlaySound(sound classes.Sound) error {
	code, err := soundJS(sound)
	if err != nil {
		return err
	}
	return browser.CallJSFunction(code)
}

func soundJS(sound classes.Sound) (string, error) {
	if sound.Name == "Default" {
		return "resetSliders()", nil
	}
	values := make([]string, len(sound.Sliders))
	for i, value := range sound.Sliders {
		if math.IsNaN(value) || math.IsInf(value, 0) {
			return "", fmt.Errorf("invalid slider %d", i)
		}
		values[i] = strconv.FormatFloat(value, 'f', -1, 64)
	}
	name, err := json.Marshal(sound.Name)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("setPreset(%s,%s);", strings.Join(values, ","), name), nil
}
