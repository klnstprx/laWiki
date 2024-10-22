# Start from the official Go image
FROM golang:1.23-alpine

# Set the Current Working Directory inside the container
WORKDIR /app

# Install git for fetching dependencies
RUN apk update && apk add --no-cache git

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the Go app
RUN go build -o api-gateway .

# Expose port (make sure it matches the one in your .env or default)
EXPOSE 8000

# Command to run the executable
CMD ["./api-gateway"]
