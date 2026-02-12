package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/HWCronicus/ssh-resume/src/server"
	logger "github.com/HWCronicus/ssh-resume/src/utils"
	"github.com/charmbracelet/lipgloss"
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
	maxHeight int
	maxWidth  int
}

func main() {

	godotenv.Load()

	config := config{
		sshPort:   os.Getenv("SSH_PORT"),
		host:      os.Getenv("HOST"),
		minHeight: 30,
		minWidth:  120,
		maxHeight: 50,
		maxWidth:  200,
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
	go server.StartWishServer(config.host, config.sshPort, config.minHeight, config.minWidth, config.maxHeight, config.maxWidth, config.renderer)
	// Uncomment the following lines to run the TUI locally instead of via SSH
	// This will start the TUI application directly in the terminal

	// time.Sleep(500 * time.Microsecond)
	// terminalWidth, terminalHeight, err := term.GetSize(os.Stdout.Fd())
	// if err != nil {
	// 	logger.LogError("TUI Unable to get terminal size", err)
	// 	os.Exit(1)
	// }
	// logger.LogInfo(fmt.Sprintf("Terminal size: %dx%d", terminalWidth, terminalHeight))
	// p := tea.NewProgram(models.InitialModel(config.minHeight, config.minWidth, config.maxHeight, config.maxWidth, terminalHeight, terminalWidth, config.renderer), tea.WithAltScreen(), tea.WithMouseCellMotion())
	// if _, err := p.Run(); err != nil {
	// 	logger.LogError("TUI application failed", err)
	// 	os.Exit(1)
	// }

	// logger.LogInfo("TUI application closed, SSH server still running")

	// Wait for interrupt signal
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	<-sigChan

	logger.LogInfo("Shutdown signal received")
	fmt.Println("\nShutting down...")
	os.Exit(0)
}
