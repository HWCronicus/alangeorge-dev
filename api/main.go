package main

import (
	"log"
	"net/http"
	"os"

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
		port:    os.Getenv("PORT"),
		sshPort: os.Getenv("SSH_PORT"),
		host:    os.Getenv("HOST"),
	}

	emailConfig = EmailConfig{
		Host:        os.Getenv("SMTP_HOST"),
		Port:        os.Getenv("SMTP_PORT"),
		Username:    os.Getenv("SMTP_USERNAME"),
		Password:    os.Getenv("SMTP_PASSWORD"),
		SenderName:  os.Getenv("SMTP_SENDER_NAME"),
		SenderEmail: os.Getenv("SMTP_SENDER_EMAIL"),
		ToEmail:     os.Getenv("SMTP_TO_EMAIL"),
	}

	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/", fs)

	http.HandleFunc("/email", SendEmailHandler)

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		handleWebSocket(w, r, config.host, config.sshPort)
	})

	if config.port == "" {
		config.port = "8080"
	}

	log.Println("Starting HTTP server on " + config.host + ":" + config.port)
	log.Fatal(http.ListenAndServe(":"+config.port, nil))
}
