FROM golang:1.23-alpine

WORKDIR /app

RUN apk update && apk add --no-cache git

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o api-gateway . && chmod +x api-gateway

EXPOSE 8000

CMD ["./api-gateway"]
