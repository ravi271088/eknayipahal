# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Build and Run Commands
- Build: `go build -o ngo-site main.go`
- Run: `go run main.go`
- Test: `go test ./...`
- Run a single test: `go test -v -run TestName main.go`
- Docker Build: `docker build -t ngo-site .`
- Docker Run: `docker run -p 8080:8080 ngo-site`

## Architecture
The project is a simple Go web application using the Gin framework to serve a static site for an NGO.

- `main.go`: The entry point of the application. It configures the Gin router, defines routes (`/`, `/about`, `/contact`), and handles template rendering.
- `static/`: Contains static assets (CSS and images).
- `templates/`: Contains HTML templates using `text/template`. All pages use `layout.html` as a base.
- `Dockerfile`: A multi-stage build that compiles the Go binary and packages it into a slim runtime image.
