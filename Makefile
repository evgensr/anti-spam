GO ?= go


.PHONY: build
build:
	$(GO) build -o ./bin/anti-spam cmd/anti-spam/main.go

.PHONY: run
run: build
	./bin/anti-spam

.PHONY: download
download:
	@echo "[*] $@"
	$(GO) mod download

.PHONY: check
check:
	@echo "[*] $@"
	staticcheck ./...
	$(GO) vet ./...

.PHONY: up
up:
	@echo "[*] $@"
	docker compose up -d



