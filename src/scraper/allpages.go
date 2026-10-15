package scraper

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/aar072/mynoise-tui/classes"
)

func FetchPresets() ([]classes.ScraperPreset, error) {
	doc, err := fetchDocument("https://mynoise.net/noiseMachines.php")
	if err != nil {
		return nil, err
	}
	baseURL, _ := url.Parse("https://mynoise.net/noiseMachines.php")

	var presets []classes.ScraperPreset
	doc.Find("a").Each(func(i int, s *goquery.Selection) {
		onmouseover, exists := s.Attr("onmouseover")
		if !exists || !strings.HasPrefix(onmouseover, "play(") {
			return
		}

		href, _ := s.Attr("href")
		text := strings.TrimSpace(s.Text())

		var category string

		// Step 1: get the closest ancestor div with class nestedSection
		ancestorDiv := s.Closest("div.nestedSection")
		if ancestorDiv.Length() > 0 {
			// Step 2: find the h1 inside that div
			h1 := ancestorDiv.Find("h1").First()
			if h1.Length() > 0 {
				category = strings.TrimSpace(h1.Text())
			}
		}

		if href != "" && text != "" {
			target, err := url.Parse(href)
			if err != nil {
				return
			}
			target = baseURL.ResolveReference(target)
			if target.Scheme != "https" || target.Host != "mynoise.net" {
				return
			}
			href = target.String()
			if category != "" {
				presets = append(presets, classes.ScraperPreset{Title: text, URL: href, Category: category})
			}
		}
	})

	if len(presets) == 0 {
		return nil, fmt.Errorf("no presets found on myNoise catalog page")
	}
	return presets, nil
}
