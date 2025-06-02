# Start with the official Golang image
FROM golang:1.23 AS builder

# Set the working directory
WORKDIR /app

# Copy the Go application source code
COPY main.go .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o secret-app main.go

# Use a minimal base image to run the compiled binary
FROM alpine:latest

# Set the working directory
WORKDIR /app

# Copy the compiled binary from the builder image
COPY --from=builder /app/secret-app .

RUN chmod +x /app/secret-app

# Expose the port that the app will run on
EXPOSE 8080

# Run the Go application
CMD ["./secret-app"]
