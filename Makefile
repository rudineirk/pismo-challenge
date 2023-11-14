install: install-linter install-mockgen install-deps

install-deps:
	go mod download
	go mod vendor

install-linter:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.55.2

install-mockgen:
	go install go.uber.org/mock/mockgen@v0.3.0

lint:
	golangci-lint run ./...

gen-mocks:
	mockgen -source ./pkg/domains/accounts/repository.go \
		-destination ./pkg/domains/accounts/mocks/repository_mock.go

test:
	go test ./...

test-unit:
	go test ./pkg/...

test-integration:
	go test ./tests/...

test-coverage:
	go test -coverpkg=./... -coverprofile=./coverage.out ./... && \
		grep -v mock coverage.out > tmpcoverage && mv tmpcoverage coverage.out
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
