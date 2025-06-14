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
	navMutex    sync.Mutex // Add a mutex for thread-safe navigation
)

func InitBrowser() error {
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", true),
		chromedp.Flag("disable-gpu", false),
		chromedp.Flag("autoplay-policy", "no-user-gesture-required"),
		chromedp.Flag("mute-audio", false),
		chromedp.Flag("no-sandbox", true),
		chromedp.Flag("remote-allow-origins", "*"), // Important for handling origins
	)

	allocCtx, allocCancel = chromedp.NewExecAllocator(context.Background(), opts...)
	ctx, ctxCancel = chromedp.NewContext(allocCtx)

	// Start browser with an empty page
	err := chromedp.Run(ctx, chromedp.Navigate("about:blank"))
	if err != nil {
		return err
	}

	return nil
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

	ctxTimeout, cancel := context.WithTimeout(ctx, 30*time.Second) // Increased timeout
	defer cancel()

	return chromedp.Run(ctxTimeout,
		chromedp.Navigate(url),
		chromedp.WaitReady(`body`, chromedp.ByQuery), // Use WaitReady instead of WaitVisible
	)
}
