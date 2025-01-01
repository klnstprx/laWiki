FROM golang:1.23-alpine AS builder

WORKDIR /app

RUN apk update && apk add --no-cache git ca-certificates

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o api-gateway . && chmod +x api-gateway

# Final stage
FROM alpine:3.18

RUN apk add --no-cache ca-certificates

WORKDIR /app

COPY --from=builder /app/api-gateway .
COPY --from=builder /app/docs/ ./docs/

EXPOSE 8000

CMD ["./api-gateway"]
