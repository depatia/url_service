run: # запустить приложение
	go run ./cmd/main.go

lint: ## запустить golangci-lint
	golangci-lint run ./...

vet: ## запустить go vet
	go vet ./...

fmt: ## форматировать код
	go fmt ./...
