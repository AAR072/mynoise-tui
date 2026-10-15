package browser_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/aar072/mynoise-tui/browser"
	"github.com/aar072/mynoise-tui/classes"
	"github.com/aar072/mynoise-tui/player"
)

func TestBrowserPlayback(t *testing.T) {
	if os.Getenv("MYNOISE_BROWSER_TESTS") != "1" {
		t.Skip("set MYNOISE_BROWSER_TESTS=1 to launch Chrome")
	}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `<script>
  window.defaultLevels = function() { return [0.42,0,0,0,0,0,0,0,0,0]; };
  window.resetSliders = function() { setPreset(...defaultLevels(), "Default"); };
  window.setPreset = function(...args) { window.lastName = args[10]; window.level = args[0]; };
  </script><body>Generator</body>`)
	}))
	defer server.Close()
	defer browser.ShutdownBrowser()
	if err := player.PlayPreset(classes.Preset{Data: classes.ScraperPreset{URL: server.URL}}); err != nil {
		t.Fatal(err)
	}
	if err := browser.CallJSFunction(`if (lastName !== "Default" || level !== 0.42) throw new Error('bad default')`); err != nil {
		t.Fatal(err)
	}
	sound := classes.Sound{Name: "A \"quote\" \\ slash\n", Sliders: [10]float64{0.123456789}}
	if err := player.PlaySound(sound); err != nil {
		t.Fatal(err)
	}
	if err := browser.CallJSFunction(`if (lastName !== "A \"quote\" \\ slash\n" || level !== 0.123456789) throw new Error('bad sound')`); err != nil {
		t.Fatal(err)
	}
	if err := browser.CallJSFunction(`throw new Error('expected playback error')`); err == nil {
		t.Fatal("JavaScript errors were swallowed")
	}
	if err := player.PlaySound(classes.DefaultSound); err != nil {
		t.Fatal(err)
	}
	if err := browser.NavigateTo(server.URL); err != nil {
		t.Fatal(err)
	}
	browser.ShutdownBrowser()
	if err := browser.NavigateTo(server.URL); err == nil {
		t.Fatal("browser restarted after shutdown")
	}
}
