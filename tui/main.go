package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/HWCronicus/ssh-resume/src/models"
	"github.com/HWCronicus/ssh-resume/src/server"
	logger "github.com/HWCronicus/ssh-resume/src/utils"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/term"
	"github.com/joho/godotenv"
	"github.com/muesli/termenv"
)

type config struct {
	port      string
	sshPort   string
	host      string
	renderer  *lipgloss.Renderer
	minHeight int
	minWidth  int
}

func main() {

	godotenv.Load()

	config := config{
		sshPort:   os.Getenv("SSH_PORT"),
		host:      os.Getenv("HOST"),
		minHeight: 50,
		minWidth:  150,
	}
	//Start logger
	if err := logger.InitLogger(); err != nil {
		fmt.Printf("Failed to initialize logger: %v\n", err)
		os.Exit(1)
	}
	defer logger.CloseLogger()

	config.renderer = lipgloss.NewRenderer(os.Stdout, termenv.WithColorCache(true))
	config.renderer.SetColorProfile(termenv.TrueColor)

	logger.LogInfo("Starting SSH Resume application")

	logger.LogInfo(fmt.Sprintf("Starting wish servers on %s:%s", config.host, config.sshPort))

	// Start both servers asynchronously
	go server.StartWishServer(config.host, config.sshPort, config.minHeight, config.minWidth, config.renderer)
	// Uncomment the following lines to run the TUI locally instead of via SSH
	// This will start the TUI application directly in the terminal

	time.Sleep(500 * time.Microsecond)
	width, height, err := term.GetSize(os.Stdout.Fd())
	if err != nil {
		width, height = config.minWidth, config.minHeight
	}
	logger.LogInfo(fmt.Sprintf("Terminal size: %dx%d", width, height))
	p := tea.NewProgram(models.InitialModel(height, width, config.renderer), tea.WithAltScreen(), tea.WithMouseCellMotion())

	if _, err := p.Run(); err != nil {
		logger.LogError("TUI application failed", err)
		os.Exit(1)
	}

	logger.LogInfo("TUI application closed, SSH server still running")

	// Wait for interrupt signal
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	<-sigChan

	logger.LogInfo("Shutdown signal received")
	fmt.Println("\nShutting down...")
	os.Exit(0)
}
