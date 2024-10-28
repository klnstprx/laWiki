FROM golang:1.23-alpine

WORKDIR /app

RUN apk update && apk add --no-cache git

COPY go.mod go.sum ./

RUN go mod tidy

COPY . .

RUN go build -o wiki-service . && chmod +x wiki-service

EXPOSE 8001

CMD ["./wiki-service"]
