run:
	LOG_FORMAT=cli go run ./cmd/main.go

build:
	docker build -t rudineirk/pismo-challenge .

install:
	go mod download
	go mod vendor

install-linter:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.55.2

lint:
	golangci-lint run ./...
