FROM golang as builder

WORKDIR /app

ENV GO111MODULE=on
ENV CGO_ENABLED=0

COPY . .

RUN make install-deps
RUN make build

FROM debian:stable-slim

WORKDIR /app

RUN apt update && apt install -y \
    ca-certificates \
    curl \
  && rm -rf /var/lib/apt/lists/*
RUN update-ca-certificates

COPY scripts/db /app/scripts/db
COPY --from=builder /app/main /app/server

EXPOSE 3000
CMD ["/app/server"]
