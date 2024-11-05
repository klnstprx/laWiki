FROM golang:1.23-bullseye

WORKDIR /app

RUN apt-get update && apt-get install -y --no-install-recommends git ca-certificates && rm -rf /var/lib/apt/lists/*

COPY go.mod go.sum ./

RUN go mod tidy

COPY . .

RUN go build -o wiki-service . && chmod +x wiki-service

EXPOSE 8001

CMD ["./wiki-service"]
