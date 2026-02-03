package main

import (
	"encoding/json"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/websocket"
	"golang.org/x/crypto/ssh"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type WSMessage struct {
	Type string `json:"type"`
	Data string `json:"data,omitempty"`
	Cols int    `json:"cols,omitempty"`
	Rows int    `json:"rows,omitempty"`
}

func UpgradeWebSocket(w http.ResponseWriter, r *http.Request) *websocket.Conn {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Upgrade error: %v", err)
		os.Exit(1)
	}
	return conn

}

func ConnectToSSH(host, port string, connection *websocket.Conn) *ssh.Client {
	hostLocation := net.JoinHostPort(host, port)
	config := &ssh.ClientConfig{
		User: "user",
		Auth: []ssh.AuthMethod{
			ssh.Password("password"),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         10 * time.Second,
	}
	log.Printf("SSH: Attempting to connect to %s", hostLocation)
	client, err := ssh.Dial("tcp", hostLocation, config)
	if err != nil {
		log.Printf("SSH dial error: %v", err)
		connection.WriteMessage(websocket.TextMessage, []byte("\r\nFailed to connect to SSH server\r\n"))
		return nil
	}

	return client
}

func CreateNewSSHSession(client *ssh.Client) *ssh.Session {
	session, err := client.NewSession()
	if err != nil {
		log.Printf("SSH session error: %v", err)
		return nil
	}
	return session
}

func handleWebSocket(w http.ResponseWriter, r *http.Request, sshHost, sshPort string) {
	conn := UpgradeWebSocket(w, r)
	defer conn.Close()

	client := ConnectToSSH(sshHost, sshPort, conn)
	if client == nil {
		return
	}
	defer client.Close()

	session := CreateNewSSHSession(client)
	defer session.Close()

	// Setup terminal modes
	modes := ssh.TerminalModes{
		ssh.ECHO:          1,
		ssh.TTY_OP_ISPEED: 14400,
		ssh.TTY_OP_OSPEED: 14400,
	}

	if err := session.RequestPty("xterm-256color", 40, 80, modes); err != nil {
		log.Printf("PTY error: %v", err)
		return
	}

	// Get stdin/stdout pipes
	sshIn, err := session.StdinPipe()
	if err != nil {
		log.Printf("Stdin pipe error: %v", err)
		return
	}

	sshOut, err := session.StdoutPipe()
	if err != nil {
		log.Printf("Stdout pipe error: %v", err)
		return
	}

	sshErr, err := session.StderrPipe()
	if err != nil {
		log.Printf("Stderr pipe error: %v", err)
		return
	}

	// Start shell
	if err := session.Shell(); err != nil {
		log.Printf("Shell error: %v", err)
		return
	}

	// Handle SSH output -> WebSocket
	go func() {
		buf := make([]byte, 32*1024)
		for {
			n, err := sshOut.Read(buf)
			if err != nil {
				if err != io.EOF {
					log.Printf("SSH read error: %v", err)
				}
				return
			}
			if n > 0 {
				if err := conn.WriteMessage(websocket.BinaryMessage, buf[:n]); err != nil {
					log.Printf("WebSocket write error: %v", err)
					return
				}
			}
		}
	}()

	// Handle SSH stderr -> WebSocket
	go func() {
		buf := make([]byte, 32*1024)
		for {
			n, err := sshErr.Read(buf)
			if err != nil {
				if err != io.EOF {
					log.Printf("SSH stderr read error: %v", err)
				}
				return
			}
			if n > 0 {
				if err := conn.WriteMessage(websocket.BinaryMessage, buf[:n]); err != nil {
					log.Printf("WebSocket write error: %v", err)
					return
				}
			}
		}
	}()

	// Handle WebSocket input -> SSH
	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("WebSocket read error: %v", err)
			}
			break
		}

		var msg WSMessage
		if err := json.Unmarshal(message, &msg); err != nil {
			log.Printf("JSON unmarshal error: %v", err)
			continue
		}

		switch msg.Type {
		case "input":
			if _, err := sshIn.Write([]byte(msg.Data)); err != nil {
				log.Printf("SSH write error: %v", err)
				return
			}
		case "resize":
			if msg.Cols > 0 && msg.Rows > 0 {
				if err := session.WindowChange(msg.Rows, msg.Cols); err != nil {
					log.Printf("Window change error: %v", err)
				}
			}
		}
	}

	session.Wait()
}
