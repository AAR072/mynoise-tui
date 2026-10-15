package scraper

import (
	"math"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/aar072/mynoise-tui/classes"
)

// FetchSounds fetches a generator page once. Default playback uses resetSliders,
// so extracting default slider values from JavaScript is unnecessary.
func FetchSounds(url string) ([]classes.Sound, error) {
	doc, err := fetchDocument(url)
	if err != nil {
		return nil, err
	}

	type soundData struct {
		Onclick string
		Name    string
	}
	var soundDatas []soundData

	// Find the h2 with text "Presets"
	doc.Find("h2").EachWithBreak(func(i int, s *goquery.Selection) bool {
		if strings.TrimSpace(s.Text()) == "Presets" {
			// Get the next sibling <p>
			p := s.NextFiltered("p")
			if p.Length() == 0 {
				return false
			}

			// Collect onclick and inner text from each span
			p.Find("span.actionlink").Each(func(i int, span *goquery.Selection) {
				if val, exists := span.Attr("onclick"); exists {
					if val != "window.location='/login.php'" {
						name := strings.TrimSpace(span.Text())
						soundDatas = append(soundDatas, soundData{
							Onclick: val,
							Name:    name,
						})
					}
				}
			})
			return false
		}
		return true
	})

	sounds := []classes.Sound{classes.DefaultSound}

	for _, data := range soundDatas {
		if !strings.HasPrefix(data.Onclick, "setPreset(") {
			continue
		}
		oc := strings.TrimPrefix(data.Onclick, "setPreset(")
		oc = strings.TrimSuffix(oc, ");")

		parts := splitIgnoringQuotes(oc, ',')

		// Only need 10 slider values, name comes from span text
		if len(parts) < 10 {
			continue
		}

		var sliders [10]float64
		valid := true
		for i := 0; i < 10; i++ {
			val, err := strconv.ParseFloat(strings.TrimSpace(parts[i]), 64)
			if err != nil || math.IsNaN(val) || math.IsInf(val, 0) {
				valid = false
				break
			}
			sliders[i] = val
		}

		if valid {
			sounds = append(sounds, classes.Sound{
				Name:    data.Name, // Use collected span text
				Sliders: sliders,
			})
		}
	}

	return sounds, nil
}

func splitIgnoringQuotes(s string, sep rune) []string {
	var parts []string
	var buf strings.Builder
	inSingleQuotes := false
	inDoubleQuotes := false

	for _, r := range s {
		switch r {
		case '\'':
			if !inDoubleQuotes {
				inSingleQuotes = !inSingleQuotes
			}
			buf.WriteRune(r)
		case '"':
			if !inSingleQuotes {
				inDoubleQuotes = !inDoubleQuotes
			}
			buf.WriteRune(r)
		case sep:
			if inSingleQuotes || inDoubleQuotes {
				buf.WriteRune(r)
			} else {
				parts = append(parts, buf.String())
				buf.Reset()
			}
		default:
			buf.WriteRune(r)
		}
	}
	if buf.Len() > 0 {
		parts = append(parts, buf.String())
	}
	return parts
}
