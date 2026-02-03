package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/HWCronicus/ssh-resume/src/server"
	logger "github.com/HWCronicus/ssh-resume/src/utils"
	"github.com/joho/godotenv"
)

type config struct {
	port    string
	sshPort string
	host    string
}

func main() {

	godotenv.Load()

	config := config{
		sshPort: os.Getenv("SSH_PORT"),
		host:    os.Getenv("HOST"),
	}
	//Start logger
	if err := logger.InitLogger(); err != nil {
		fmt.Printf("Failed to initialize logger: %v\n", err)
		os.Exit(1)
	}
	defer logger.CloseLogger()

	logger.LogInfo("Starting SSH Resume application")

	logger.LogInfo(fmt.Sprintf("Starting wish servers on %s:%s", config.host, config.sshPort))

	// Start both servers asynchronously
	go server.StartWishServer(config.host, config.sshPort)
	// Uncomment the following lines to run the TUI locally instead of via SSH
	// This will start the TUI application directly in the terminal

	/*
		Give servers a moment to start
		time.Sleep(500 * time.Microsecond)
		width, height, err := term.GetSize(os.Stdout.Fd())
		if err != nil {
			width, height = 200, 50
		}
		logger.LogInfo(fmt.Sprintf("Terminal size: %dx%d", width, height))
		p := tea.NewProgram(models.InitialModel(height, width), tea.WithAltScreen(), tea.WithMouseCellMotion())

		if _, err := p.Run(); err != nil {
			logger.LogError("TUI application failed", err)
			os.Exit(1)
		}

		After TUI exits, keep the SSH server running
		logger.LogInfo("TUI application closed, SSH server still running")
	*/

	// Wait for interrupt signal
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	<-sigChan

	logger.LogInfo("Shutdown signal received")
	fmt.Println("\nShutting down...")
	os.Exit(0)
}
