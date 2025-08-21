FROM golang:1.24.5-bookworm

RUN apt-get update -y && apt-get install -y --no-install-recommends \
    postgresql-client && rm -rf /var/lib/apt/lists/*

RUN go install github.com/swaggo/swag/cmd/swag@latest

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN swag init && swag fmt
RUN go build -o /app/main .
CMD ["/app/main"]
