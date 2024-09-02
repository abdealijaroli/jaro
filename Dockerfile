FROM golang:1.20 AS builder

WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./
RUN go mod download

# Copy source files
COPY . .

# Build the Go binary
RUN go build -o jaro

# Create a minimal image for running the Go binary
FROM alpine:latest

WORKDIR /root/

# Copy the Go binary from the builder stage
COPY --from=builder /app/jaro .

# Copy the web files into the image
COPY web /root/web

# Expose port 8008
EXPOSE 8008

# Run the Go binary
CMD ["./jaro"]
