package scraper

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type Preset struct {
	Title    string
	URL      string
	Category string
}

func FetchPresets() ([]Preset, error) {
	res, err := http.Get("https://mynoise.net/noiseMachines.php")
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return nil, fmt.Errorf("status code error: %d %s", res.StatusCode, res.Status)
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return nil, err
	}

	var presets []Preset
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
			if strings.HasPrefix(href, "/") {
				href = "https://mynoise.net" + href
			}
			if category != "" {
				presets = append(presets, Preset{Title: text, URL: href, Category: category})
			}
		}
	})

	return presets, nil
}
