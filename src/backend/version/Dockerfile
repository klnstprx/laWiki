FROM golang:1.23-alpine AS builder

WORKDIR /app

RUN apk update && apk add --no-cache git ca-certificates

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o version-service . && chmod +x version-service

# Final stage
FROM alpine:3.18

RUN apk add --no-cache ca-certificates

WORKDIR /app

COPY --from=builder /app/version-service .

EXPOSE 8005

CMD ["./version-service"]
