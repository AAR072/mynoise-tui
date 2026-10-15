package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/aar072/mynoise-tui/browser"
	"github.com/aar072/mynoise-tui/model"
	"github.com/aar072/mynoise-tui/prefs"
	tea "github.com/charmbracelet/bubbletea"
)

func run() error {
	if err := prefs.InitConfig(); err != nil {
		return err
	}
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	defer browser.ShutdownBrowser()
	p := tea.NewProgram(model.NewModel(), tea.WithAltScreen(), tea.WithContext(ctx))
	_, runErr := p.Run()
	if errors.Is(runErr, tea.ErrProgramKilled) && ctx.Err() != nil {
		runErr = nil
	}
	return errors.Join(runErr, prefs.SavePreferences())
}

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
