# Stage 1: Build the Go binary
FROM golang:latest

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy the Go Modules manifests
COPY . .

# Download all Go modules. This will cache the result, so it's only done when the go.mod or go.sum changes
RUN go mod download -x

# Build the Go application
RUN go build -o ./main

EXPOSE 8080

CMD ["./main"]
