package scraper

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

type transportFunc func(*http.Request) (*http.Response, error)

func (f transportFunc) RoundTrip(r *http.Request) (*http.Response, error) { return f(r) }

func TestFetchSoundsFailure(t *testing.T) {
	previous := httpClient
	t.Cleanup(func() { httpClient = previous })
	httpClient = &http.Client{Transport: transportFunc(func(*http.Request) (*http.Response, error) {
		return nil, errors.New("offline")
	})}
	if _, err := FetchSounds("https://example.invalid"); err == nil {
		t.Fatal("expected network error")
	}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(503) }))
	defer server.Close()
	httpClient = server.Client()
	if _, err := FetchSounds(server.URL); err == nil {
		t.Fatal("expected HTTP status error")
	}
}

func TestFetchSoundsTimeout(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { <-r.Context().Done() }))
	defer server.Close()
	previous := httpClient
	t.Cleanup(func() { httpClient = previous })
	httpClient = &http.Client{Timeout: 30 * time.Millisecond}
	if _, err := FetchSounds(server.URL); err == nil {
		t.Fatal("expected timeout")
	}
}

func TestFetchSoundsUsesOneRequest(t *testing.T) {
	requests := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requests++
		fmt.Fprint(w, `<script>
 function defaultLevels() { return CFG.defaultCurve || [0.5,0.5,0.5,0.5,0.5,0.5,0.5,0.5,0.5,0.5]; }
 function resetSliders() { setPreset(...defaultLevels(), "Back to Default"); }
 </script><h2>Presets</h2><p>
  <span class="actionlink" onclick="setPreset(0.12345,1,2,3,4,5,6,7,8,9,this.innerHTML);">A &quot;quoted&quot; sound</span>
  <span class="actionlink" onclick="window.location='/login.php'">Locked</span>
  <span class="actionlink" onclick="setPreset(NaN,1,2,3,4,5,6,7,8,9,'bad');">Invalid</span></p>`)
	}))
	defer server.Close()
	sounds, err := FetchSounds(server.URL)
	if err != nil {
		t.Fatal(err)
	}
	if requests != 1 || len(sounds) != 2 {
		t.Fatalf("requests=%d sounds=%v", requests, sounds)
	}
	if sounds[0].Name != "Default" || sounds[1].Name != `A "quoted" sound` || sounds[1].Sliders[0] != 0.12345 {
		t.Fatalf("incorrect parsing: %+v", sounds)
	}
}

func TestCatalogResolvesURLsAndRejectsEmptyPages(t *testing.T) {
	previous := httpClient
	t.Cleanup(func() { httpClient = previous })
	html := `<div class="nestedSection"><h1>Nature</h1><a onmouseover="play('X')" href="NoiseMachines/rain.php">Rain</a></div>`
	httpClient = &http.Client{Transport: transportFunc(func(*http.Request) (*http.Response, error) {
		return &http.Response{StatusCode: 200, Body: http.NoBody}, nil
	})}
	if _, err := FetchPresets(); err == nil {
		t.Fatal("empty catalog must fail")
	}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { fmt.Fprint(w, html) }))
	defer server.Close()
	httpClient = &http.Client{Transport: transportFunc(func(r *http.Request) (*http.Response, error) {
		target := *r.URL
		replacement, _ := http.NewRequest("GET", server.URL, nil)
		r.URL = replacement.URL
		r.Host = replacement.Host
		response, err := http.DefaultTransport.RoundTrip(r)
		if response != nil {
			response.Request.URL = &target
		}
		return response, err
	})}
	presets, err := FetchPresets()
	if err != nil {
		t.Fatal(err)
	}
	if len(presets) != 1 || !strings.HasPrefix(presets[0].URL, "https://mynoise.net/NoiseMachines/") {
		t.Fatalf("bad catalog: %+v", presets)
	}
}
