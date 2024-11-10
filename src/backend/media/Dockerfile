FROM golang:1.23-alpine3.18 AS builder

WORKDIR /app

RUN apk update && apk add --no-cache git ca-certificates

COPY go.mod go.sum ./

RUN go mod tidy

COPY . .

RUN go build -o media-service . && chmod +x media-service

# Final stage
FROM alpine:3.18

RUN apk add --no-cache ca-certificates

WORKDIR /app

COPY --from=builder /app/version-service .

EXPOSE 8081

CMD ["./media-service"]
