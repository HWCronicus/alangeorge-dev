package server

import (
	"context"
	"errors"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/HWCronicus/ssh-resume/src/models"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
	"github.com/charmbracelet/ssh"
	"github.com/charmbracelet/wish"
	"github.com/charmbracelet/wish/activeterm"
	"github.com/charmbracelet/wish/bubbletea"
	"github.com/charmbracelet/wish/logging"
)

func StartWishServer(host, port string, minHeight, minWidth int, renderer *lipgloss.Renderer) {
	s, err := wish.NewServer(
		wish.WithAddress(net.JoinHostPort(host, port)),
		wish.WithHostKeyPath(".ssh/id_ed25519"),
		wish.WithMiddleware(
			bubbletea.Middleware(func(s ssh.Session) (tea.Model, []tea.ProgramOption) {
				return teaHandler(minHeight, minWidth, s, renderer)
			}),
			activeterm.Middleware(),
			logging.Middleware(),
		),
	)
	if err != nil {
		log.Error("Could not create server", "error", err)
		return
	}

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	log.Info("Starting SSH server", "host", host, "port", port)

	go func() {
		if err = s.ListenAndServe(); err != nil && !errors.Is(err, ssh.ErrServerClosed) {
			log.Error("Could not start server", "error", err)
			done <- nil
		}
	}()

	<-done
	log.Info("Stopping SSH server")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if err := s.Shutdown(ctx); err != nil && !errors.Is(err, ssh.ErrServerClosed) {
		log.Error("Could not stop server", "error", err)
	}
}

func teaHandler(minHeight, minWidth int, s ssh.Session, renderer *lipgloss.Renderer) (tea.Model, []tea.ProgramOption) {
	pty, _, ok := s.Pty()
	if !ok {
		log.Error("No PTY requested")
		return nil, nil
	}

	log.Info("New SSH session",
		"user", s.User(),
		"width", pty.Window.Width,
		"height", pty.Window.Height,
	)

	terminalWidth := pty.Window.Width
	terminalHeight := pty.Window.Height
	if terminalWidth <= 0 {
		terminalWidth = minWidth
	}
	if terminalHeight <= 0 {
		terminalHeight = minHeight
	}

	m := models.InitialModel(minHeight, minWidth, terminalHeight, terminalWidth, renderer)
	return m, []tea.ProgramOption{tea.WithAltScreen(), tea.WithMouseCellMotion()}
}
