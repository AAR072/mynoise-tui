package scraper

import (
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/aar072/mynoise-tui/classes"
	"github.com/aar072/mynoise-tui/logger"
)

func FetchPresetOnclicks(url string) []classes.Sound {
	res, _ := http.Get(url)
	defer res.Body.Close()

	doc, _ := goquery.NewDocumentFromReader(res.Body)

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

	var sounds []classes.Sound

	for _, data := range soundDatas {
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
			if err != nil {
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

	return sounds
}

func GetDefaultSound(url string) classes.Sound {
	res, err := http.Get(url)
	if err != nil {
		logger.Error("failed to fetch page: " + err.Error())
		return classes.Sound{}
	}
	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		logger.Error("failed to parse HTML: " + err.Error())
		return classes.Sound{}
	}

	var sound classes.Sound
	found := false

	doc.Find("script").EachWithBreak(func(i int, s *goquery.Selection) bool {
		script := s.Text()
		if strings.Contains(script, "function resetSliders()") {
			// Extract resetSliders function body
			reFunc := regexp.MustCompile(`function resetSliders\s*\(\s*\)\s*\{([\s\S]*?)\}`)
			matchesFunc := reFunc.FindStringSubmatch(script)
			if len(matchesFunc) < 2 {
				logger.Error("failed to extract resetSliders function body")
				return true
			}

			resetBody := matchesFunc[1]

			// Extract setPreset call inside resetSliders
			reSetPreset := regexp.MustCompile(`setPreset\s*\((.*?)\)`)
			setPresetMatch := reSetPreset.FindStringSubmatch(resetBody)
			if len(setPresetMatch) < 2 {
				logger.Error("setPreset call not found inside resetSliders")
				return true
			}

			argsStr := setPresetMatch[1]
			logger.Info("Extracted args string: " + argsStr)

			args := splitIgnoringQuotes(argsStr, ',')
			logger.Info(fmt.Sprintf("Split args count: %d", len(args)))
			for i, a := range args {
				logger.Info(fmt.Sprintf("Arg[%d]: %s", i, a))
			}

			if len(args) < 10 {
				logger.Error("expected at least 10 args for sliders, got " + strconv.Itoa(len(args)))
				return true
			}

			var sliders [10]float64
			for i := 0; i < 10; i++ {
				val, err := strconv.ParseFloat(strings.TrimSpace(args[i]), 64)
				if err != nil {
					logger.Error("invalid float in default sound: " + err.Error())
					return true
				}
				sliders[i] = val
			}

			sound = classes.Sound{
				Name:    "Default", // hardcoded name
				Sliders: sliders,
			}
			found = true
			return false // break EachWithBreak
		}
		return true
	})

	if !found {
		logger.Error("default setPreset(...) not found in page")
	}

	return sound
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
