package browser

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/chromedp/chromedp"
)

var (
	navMutex    sync.Mutex
	lifecycle   sync.Mutex
	ctx         context.Context
	ctxCancel   context.CancelFunc
	allocCancel context.CancelFunc
	closed      bool
)

// initBrowser is called with navMutex held, and starts Chrome only when needed.
func initBrowser() error {
	lifecycle.Lock()
	if closed {
		lifecycle.Unlock()
		return fmt.Errorf("browser is shut down")
	}
	if ctx != nil {
		lifecycle.Unlock()
		return nil
	}
	opts := append([]chromedp.ExecAllocatorOption{}, chromedp.DefaultExecAllocatorOptions[:]...)
	opts = append(opts,
		chromedp.Flag("autoplay-policy", "no-user-gesture-required"),
		chromedp.Flag("mute-audio", false),
	)
	allocCtx, cancelAlloc := chromedp.NewExecAllocator(context.Background(), opts...)
	browserCtx, cancelCtx := chromedp.NewContext(allocCtx)
	ctx, ctxCancel, allocCancel = browserCtx, cancelCtx, cancelAlloc
	lifecycle.Unlock()
	// The first Run must use the persistent context: cancelling a child
	// after allocation also terminates the browser it created.
	timer := time.AfterFunc(30*time.Second, cancelCtx)
	defer timer.Stop()
	if err := chromedp.Run(browserCtx); err != nil {
		cancelCtx()
		cancelAlloc()
		lifecycle.Lock()
		ctx = nil
		lifecycle.Unlock()
		return err
	}
	return nil
}

func InitBrowser() error {
	navMutex.Lock()
	defer navMutex.Unlock()
	return initBrowser()
}

func ShutdownBrowser() {
	lifecycle.Lock()
	closed = true
	cancelCtx, cancelAlloc := ctxCancel, allocCancel
	lifecycle.Unlock()
	if cancelCtx != nil {
		cancelCtx()
	}
	if cancelAlloc != nil {
		cancelAlloc()
	}
}

func NavigateTo(url string) error {
	navMutex.Lock()
	defer navMutex.Unlock()
	if err := initBrowser(); err != nil {
		return err
	}
	timeout, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()
	return chromedp.Run(timeout,
		chromedp.Navigate(url),
		chromedp.WaitReady(`body`, chromedp.ByQuery),
	)
}

func CallJSFunction(code string) error {
	navMutex.Lock()
	defer navMutex.Unlock()
	if err := initBrowser(); err != nil {
		return err
	}
	timeout, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()
	return chromedp.Run(timeout,
		chromedp.Poll(`typeof resetSliders === 'function' && typeof setPreset === 'function'`, nil),
		chromedp.Evaluate(code, nil),
	)
}
