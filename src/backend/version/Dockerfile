FROM golang:1.23-bullseye

WORKDIR /app

RUN apt-get update && apt-get install -y --no-install-recommends git ca-certificates && rm -rf /var/lib/apt/lists/*

COPY go.mod go.sum ./

RUN go mod tidy

COPY . .

RUN go build -o version-service . && chmod +x version-service

EXPOSE 8005

CMD ["./version-service"]
