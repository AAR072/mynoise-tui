package browser

import (
	"context"
	"sync"
	"time"

	"github.com/chromedp/chromedp"
)

var (
	allocCtx    context.Context
	allocCancel context.CancelFunc
	ctx         context.Context
	ctxCancel   context.CancelFunc
	navMutex    sync.Mutex
	statusMutex sync.Mutex
	lastStatus  bool
	lastCheck   time.Time
)

func InitBrowser() error {
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", true),
		chromedp.Flag("disable-gpu", false),
		chromedp.Flag("autoplay-policy", "no-user-gesture-required"),
		chromedp.Flag("mute-audio", false),
		chromedp.Flag("no-sandbox", true),
		chromedp.Flag("remote-allow-origins", "*"),
	)

	allocCtx, allocCancel = chromedp.NewExecAllocator(context.Background(), opts...)
	ctx, ctxCancel = chromedp.NewContext(allocCtx)

	return chromedp.Run(ctx, chromedp.Navigate("about:blank"))
}

func ShutdownBrowser() {
	if ctxCancel != nil {
		ctxCancel()
	}
	if allocCancel != nil {
		allocCancel()
	}
}

func NavigateTo(url string) error {
	navMutex.Lock()
	defer navMutex.Unlock()

	ctxTimeout, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	return chromedp.Run(ctxTimeout,
		chromedp.Navigate(url),
		chromedp.WaitReady(`body`, chromedp.ByQuery),
	)
}

func IsLoading() bool {
	statusMutex.Lock()
	defer statusMutex.Unlock()

	// Rate limit checks to max once per 500ms
	if time.Since(lastCheck) < 500*time.Millisecond {
		return false
	}

	var status bool
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	chromedp.Run(ctx,
		chromedp.Evaluate(`
			(() => {
				const el = document.querySelector("div.msg#msg");
				if (!el) return "absent";
				return el.textContent == "Now Playing...";
			})()
		`, &status),
	)

	lastStatus = status
	lastCheck = time.Now()
	return !status
}

func CallJSFunction(jsCode string) (string, error) {
	navMutex.Lock()
	defer navMutex.Unlock()

	ctxTimeout, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	var result string
	err := chromedp.Run(ctxTimeout,
		chromedp.Evaluate(jsCode, &result),
	)
	return result, err
}
