package scraper

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/aar072/mynoise-tui/classes"
)

func FetchPresetOnclicks(url string) []classes.Sound {
	res, _ := http.Get(url)
	defer res.Body.Close()

	doc, _ := goquery.NewDocumentFromReader(res.Body)

	var onclicks []string

	// Find the h2 with text "Presets"
	doc.Find("h2").EachWithBreak(func(i int, s *goquery.Selection) bool {
		if strings.TrimSpace(s.Text()) == "Presets" {
			// Get the next sibling <p>
			p := s.NextFiltered("p")
			if p.Length() == 0 {
				return false // Stop if no <p> sibling
			}

			// Find all <span class="actionlink"> inside that <p>
			p.Find("span.actionlink").Each(func(i int, span *goquery.Selection) {
				if val, exists := span.Attr("onclick"); exists {
					if val != "window.location='/login.php'" {
						onclicks = append(onclicks, val)
					}
				}
			})

			return false // We found the target <h2>, so stop
		}
		return true // keep looking
	})

	var sounds []classes.Sound

	for _, oc := range onclicks {
		oc = strings.TrimPrefix(oc, "setPreset(")
		oc = strings.TrimSuffix(oc, ");")

		parts := splitIgnoringQuotes(oc, ',')

		if len(parts) != 11 {
			continue // Expect 10 sliders + 1 name
		}

		var sliders [10]float64
		valid := true
		// Parse sliders first (parts[0] to parts[9])
		for i := 0; i < 10; i++ {
			val, err := strconv.ParseFloat(strings.TrimSpace(parts[i]), 64)
			if err != nil {
				valid = false
				break
			}
			sliders[i] = val
		}

		// Name is last part (parts[10])
		name := strings.Trim(parts[10], " '\"")

		if valid {
			sounds = append(sounds, classes.Sound{
				Name:    name,
				Sliders: sliders,
			})
		}
	}

	return sounds
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
