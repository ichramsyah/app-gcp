FROM golang:alpine AS builder

WORKDIR /app

# Initialize go.mod inside the container since Go isn't installed on the host
RUN go mod init github.com/ichramsyah/app-gcp-backend

# Copy the Go source files into the container
COPY . .

# Generate go.sum and download all dependencies based on the source code imports
RUN go mod tidy

# Build the application
RUN go build -o main .

# Final stage
FROM alpine:latest

WORKDIR /app

# Copy the built binary from the builder stage
COPY --from=builder /app/main .

EXPOSE 8080

CMD ["./main"]
