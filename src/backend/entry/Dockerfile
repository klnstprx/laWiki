FROM golang:1.23-bullseye

WORKDIR /app

RUN apt-get update && apt-get install -y --no-install-recommends git ca-certificates && rm -rf /var/lib/apt/lists/*

COPY go.mod go.sum ./

RUN go mod tidy

COPY . .

RUN go build -o entry-service . && chmod +x entry-service

EXPOSE 8002

CMD ["./entry-service"]
