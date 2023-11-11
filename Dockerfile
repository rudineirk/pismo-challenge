FROM golang as builder

WORKDIR /app

ENV GO111MODULE=on
ENV CGO_ENABLED=0

RUN go install github.com/rubenv/sql-migrate/sql-migrate@latest

COPY . .

RUN make install-deps
RUN make build

FROM debian:stable-slim

RUN apt update && apt install -y \
    ca-certificates \
    curl \
  && rm -rf /var/lib/apt/lists/*
RUN update-ca-certificates

COPY --from=builder /go/bin/sql-migrate /bin/sql-migrate
COPY scripts/db /etc/migrate
COPY scripts/run.sh /run.sh
COPY --from=builder /app/main /server

EXPOSE 3000
CMD ["/run.sh"]
