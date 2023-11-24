# pismo-challenge

[![CI](https://github.com/rudineirk/pismo-challenge/actions/workflows/ci.yml/badge.svg)](https://github.com/rudineirk/pismo-challenge/actions/workflows/ci.yml)

This is a project to manage cardholder accounts and transactions. It was implemented using Golang, based on the clean architecture pattern.

## API documentation 📖

The API documentation is available as an OpenAPI yaml [here](./docs/openapi.yaml). You can also view the rendered version [on this link](http://rudineirk.github.io/pismo-challenge/api-docs/)

## Project structure 🏗️

```sh
docs/                 # API documentation
cmd/
  main.go             # main function, setups everything and starts the server
pkg/
  domains/            # app business rules domains
    accounts/
      api.go          # HTTP API handlers
      entity.go       # entity of the domain
      repository.go   # SQL database repository
      service.go      # service responsible for the business rules/use cases
      service_test.go # service/use cases unit tests
    operationtypes/
    transactions/
  infra/              # infrastructure required to run the project
    config/           # env vars config, to be loaded with k8s secrets or some tool like this
    database/         # PostgreSQL database setup tools
    httprouter/       # Gin HTTP router setup
    logger/           # zerolog structured (json) logger
    signalhandler/    # shutdown signals handler, to allow zero downtime restarts/upgrades
  utils/              # helpers/tools used accross the project
scripts/
  db/
    migrations/       # database migrations (using sql-migrate)
tests/                # integration tests
  accounts/
    accounts_test.go  # accounts APIs integration tests
```

## How to run 🚀

Requirements:
* docker
* docker-compose (or use the new builtin `docker compose` commmand)

To run the project, just run it using `docker-compose`, it will setup everything necessary to execute the service:

```sh
docker-compose up -d
```

After it starts everything, you can call the service APIs on the address `http://localhost:3000`:

```sh
curl -v -X POST \
  -H 'Content-Type: application/json' \
  http://localhost:3000/accounts \
  -d '{"document_number":"91219245000160"}'

curl -v http://localhost:3000/accounts/1

curl -v -X POST \
  -H 'Content-Type: application/json' \
  http://localhost:3000/transactions \
  -d '{"account_id":1,"operation_type_id":1,"amount":-1.25}'
```

## Tests 🧑‍💻

The tests are being run in the Github Actions CI of the repository, but if you wish to run it locally,
first start the database using `docker-compose`:

```sh
docker-compose up -d postgres
```

Then run the linters and tests:

```sh
make install
make lint
make test

# to run only unit tests, use this command:
make test-unit

# to run only integration tests, use this command:
make test-integration

# to view the project test coverage, use this command:
make test-coverage
```

You can also view the coverage report [on this link](http://rudineirk.github.io/pismo-challenge/coverage/)

## Improvements for this project 📈

To run this project in production, there are some things that could be implemented before to ensure it runs smoothly:

* Authentication / Authorization
  * This service should be run behind an API Gateway that provides authentication / authorization rules (like Kong, Traefik)
  * Or maybe it could use a well known identity service, like Keycloak, even a SaaS one like Auth0
* Error tracking
  * using Sentry or some tool like this to track errors
* Tracing
  * OpenTracing is the most flexible tool for tracing, you can use it with multiple providers (Zipkin, Jaeger, NewRelic)
* Monitoring
  * The service could export some usage metrics to Prometheus and make a dashboard on Grafana to monitor relevant business data in real time
