# Build frontend
FROM oven/bun:1-alpine AS frontend-builder

WORKDIR /app
COPY app/package*.json ./
COPY app/bun.lockb* ./
RUN bun install
COPY app/ .
RUN bun run build

# Build Go applications
FROM golang:1.25.4-alpine AS go-builder

WORKDIR /app
COPY api/go.mod api/go.sum* ./api/
RUN cd api && go mod download
COPY tui/go.mod tui/go.sum* ./tui/
RUN cd tui && go mod download
COPY api/ ./api/
COPY tui/ ./tui/

RUN cd api && CGO_ENABLED=0 GOOS=linux go build -o server .
RUN cd tui && CGO_ENABLED=0 GOOS=linux go build -o tui .

# Final production image
FROM alpine:3.19
RUN apk --no-cache add ca-certificates ncurses-terminfo-base openssh-keygen

WORKDIR /app

ENV TERM=xterm-256color
ENV COLORTERM=truecolor

# Copy Go binaries
COPY --from=go-builder /app/api/server .
COPY --from=go-builder /app/tui/tui .

# Copy TUI runtime assets
COPY --from=go-builder /app/tui/assets/ ./assets/
COPY --from=go-builder /app/tui/content/ ./content/

# Copy frontend static files
COPY --from=frontend-builder /app/dist ./static

# Host key is created/loaded at runtime by the entrypoint script
COPY entrypoint.sh .
RUN chmod +x entrypoint.sh && mkdir -p .ssh

EXPOSE 8282 42069
CMD ["./entrypoint.sh"]