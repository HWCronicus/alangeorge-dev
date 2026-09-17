# alangeorge-dev

Personal portfolio site with a web frontend and an SSH terminal interface.

## Project parts

- `app/` - React Vite frontend with React Three Fiber 3D effects .
- `api/` - Go HTTP API. It serves the site, handles contact email, and connects the web terminal to SSH.
- `tui/` - Go terminal UI and SSH server for the interactive resume.

## Run locally

Requirements: Docker and Docker Compose.

1. Create a `.env` file in the repository root if your local setup needs environment variables.
2. Start the development services:

```sh
npm run dev
```

Open these addresses while the services are running:

- Frontend: http://localhost:5173
- API: http://localhost:8282
- SSH terminal: `ssh -p 42069 localhost`

The frontend, API, and TUI are mounted into their containers, so code changes reload during development.

## Common commands

```sh
npm run dev:build  # build and start development containers
npm run dev:down   # stop development containers
npm run logs       # view production container logs
npm run clean      # stop containers and remove volumes
```

## Production

Build and start the production container with:

```sh
npm run up:build
```

Stop it with:

```sh
npm run down
```

