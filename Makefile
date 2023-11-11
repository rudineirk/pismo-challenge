install: install-linter install-sql-migrate install-deps

install-deps:
	go mod download
	go mod vendor

install-linter:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.55.2

install-sql-migrate:
	go install github.com/rubenv/sql-migrate/sql-migrate@latest

lint:
	golangci-lint run ./...

run:
	LOG_FORMAT=cli go run ./cmd/main.go

build:
	go build -o ./main ./cmd/main.go

build-docker:
	docker build -t rudineirk/pismo-challenge .

migrate:
	cd scripts/db && \
		DATABASE_URL="postgresql://dev:development@127.0.0.1/pismo_challenge?sslmode=disable" \
		sql-migrate up
