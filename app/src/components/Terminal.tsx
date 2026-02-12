import { useEffect, useRef, useState } from "react";
import { Terminal as XTerm } from "xterm";
import { FitAddon } from "xterm-addon-fit";
import { WebLinksAddon } from "xterm-addon-web-links";
import "xterm/css/xterm.css";
import "../styles/terminal.css";

type Status = "connecting" | "connected" | "disconnected";

export default function Terminal() {
  const terminalRef = useRef<HTMLDivElement>(null);
  const [status, setStatus] = useState<Status>("connecting");

  useEffect(() => {
    if (!terminalRef.current) return;

    const term = new XTerm({
      cursorBlink: true,
      fontSize: 14,
      fontFamily: '"Cascadia Code", "Fira Code", monospace',
      theme: {
        background: "#1e1e1e",
        foreground: "#d4d4d4",
        cursor: "#d4d4d4",
        black: "#000000",
        red: "#cd3131",
        green: "#0dbc79",
        yellow: "#e5e510",
        blue: "#2472c8",
        magenta: "#bc3fbc",
        cyan: "#11a8cd",
        white: "#e5e5e5",
        brightBlack: "#666666",
        brightRed: "#f14c4c",
        brightGreen: "#23d18b",
        brightYellow: "#f5f543",
        brightBlue: "#3b8eea",
        brightMagenta: "#d670d6",
        brightCyan: "#29b8db",
        brightWhite: "#e5e5e5",
      },
      cols: 80,
      rows: 40,
    });

    const fitAddon = new FitAddon();
    const webLinksAddon = new WebLinksAddon();

    term.loadAddon(fitAddon);
    term.loadAddon(webLinksAddon);
    term.open(terminalRef.current);

    // Delay fit to ensure container is sized
    setTimeout(() => fitAddon.fit(), 50);

    // WebSocket connection
    const protocol = window.location.protocol === "https:" ? "wss:" : "ws:";
    let wsUrl = `${protocol}//${window.location.host}/ws`;
    if (import.meta.env.DEV) {
      wsUrl = "ws://localhost:8282/ws";
    }
    console.log("Connecting to WebSocket at", wsUrl);
    let ws: WebSocket;
    let reconnectTimeout: number;

    const connect = () => {
      setStatus("connecting");
      ws = new WebSocket(wsUrl);
      ws.binaryType = "arraybuffer";

      ws.onopen = () => {
        setStatus("connected");
        ws.send(
          JSON.stringify({
            type: "resize",
            cols: term.cols,
            rows: term.rows,
          }),
        );
      };

      ws.onmessage = (event) => {
        if (typeof event.data === "string") {
          term.write(event.data);
        } else {
          const uint8Array = new Uint8Array(event.data);
          term.write(uint8Array);
        }
      };

      ws.onerror = () => {
        setStatus("disconnected");
      };

      ws.onclose = () => {
        setStatus("disconnected");
        term.write(
          "\r\n\r\n\x1b[31mConnection closed. Reconnecting in 3 seconds...\x1b[0m\r\n",
        );
        reconnectTimeout = window.setTimeout(connect, 3000);
      };
    };

    term.onData((data) => {
      if (ws && ws.readyState === WebSocket.OPEN) {
        ws.send(JSON.stringify({ type: "input", data }));
      }
    });

    term.onResize(({ cols, rows }) => {
      if (ws && ws.readyState === WebSocket.OPEN) {
        ws.send(JSON.stringify({ type: "resize", cols, rows }));
      }
    });

    const handleResize = () => fitAddon.fit();
    window.addEventListener("resize", handleResize);

    connect();
    term.focus();

    return () => {
      clearTimeout(reconnectTimeout);
      window.removeEventListener("resize", handleResize);
      ws?.close();
      term.dispose();
    };
  }, []);

  return (
    <div className="terminal-wrapper">
      <div className="terminal-header">
        <h1>
          SSH Terminal
          <span className={`status ${status}`}>
            {status.charAt(0).toUpperCase() + status.slice(1)}
          </span>
        </h1>
      </div>
      <div className="terminal-container">
        <div ref={terminalRef} id="terminal" />
      </div>
    </div>
  );
}
