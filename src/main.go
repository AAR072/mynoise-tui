package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/aar072/mynoise-tui/browser"
	"github.com/aar072/mynoise-tui/model"
	"github.com/aar072/mynoise-tui/prefs"
	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	// Setup browser
	err := browser.InitBrowser()
	if err != nil {
		log.Fatal(err)
	}
	defer browser.ShutdownBrowser() // Ensure browser closes on exit

	// Handle OS signals for proper cleanup
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sigCh
		browser.ShutdownBrowser()
		os.Exit(0)
	}()

	// Start TUI
	_ = prefs.InitConfig()
	p := tea.NewProgram(model.NewModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
	prefs.SavePreferences()
}
