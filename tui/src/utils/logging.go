package utils

import (
	"fmt"
	"io"
	"log"
	"os"
	"sync"
	"time"
)

var (
	logFile *os.File
	logMu   sync.Mutex
)

func InitLogger() error {
	logMu.Lock()
	defer logMu.Unlock()

	if err := moveAppLogContentsToArchiveLog(); err != nil {
		return err
	}

	var err error
	logFile, err = os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}

	mw := io.MultiWriter(os.Stdout, logFile)
	log.SetOutput(mw)
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	return nil
}

func moveAppLogContentsToArchiveLog() error {
	if _, err := os.Stat("app.log"); err == nil {
		content, err := os.ReadFile("app.log")
		if err != nil {
			return err
		}

		archiveFile, err := os.OpenFile("app_archive.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
		if err != nil {
			return err
		}
		defer archiveFile.Close()

		header := fmt.Sprintf("\n=== Application start %s ===\n", time.Now().Format("2006-01-02 15:04:05"))
		if _, err := archiveFile.WriteString(header); err != nil {
			return err
		}

		if _, err := archiveFile.Write(content); err != nil {
			return err
		}
	}

	return nil
}

func CloseLogger() error {
	logMu.Lock()
	defer logMu.Unlock()

	if logFile != nil {
		return logFile.Close()
	}
	return nil
}

func LogInfo(msg string) {
	log.Printf("[INFO] %s\n", msg)
}

func LogError(msg string, err error) {
	log.Printf("[ERROR] %s: %v\n", msg, err)
}

func LogDebug(msg string) {
	log.Printf("[DEBUG] %s\n", msg)
}
