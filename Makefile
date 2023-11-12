install: install-linter install-deps

install-deps:
	go mod download
	go mod vendor

install-linter:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.55.2

lint:
	golangci-lint run ./...

test:
	go test ./...

test-unit:
	go test ./pkg/...

test-integration:
	go test ./tests/...

test-coverage:
	go test -coverpkg=./... -coverprofile=./coverage.out ./...
	go tool cover -html=./coverage.out

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
