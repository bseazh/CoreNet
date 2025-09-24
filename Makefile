SHELL := /bin/bash

.PHONY: all dev up down gen fmt tidy

all: gen

gen:
	@echo "Generate services from .api (install goctl first)"
	@echo "  goctl api go -api internal/api/upload.api -dir ."
	@echo "  goctl api go -api internal/api/filemeta.api -dir ."
	@echo "  goctl api go -api internal/api/preview.api -dir ."
	@echo "  goctl api go -api internal/api/search.api -dir ."
	@echo "  goctl api go -api internal/api/taskcenter.api -dir ."

fmt:
	go fmt ./...

tidy:
	go mod tidy

dev:
	@echo "Start services manually after goctl generation."
	@echo "Example: go run ./internal/upload.go (depending on goctl output structure)"

up:
	cd deploy && docker compose up -d

down:
	cd deploy && docker compose down
