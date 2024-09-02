# Stage 1: Build the Go binary
FROM golang:1.22.3 AS builder

WORKDIR /app

LABEL maintainer="Abdeali Jaroli <abdeali@jaro.li>"

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o jaro .

# Stage 2: Create a lightweight image to run the binary
FROM alpine:latest

WORKDIR /root/

COPY --from=builder /app/jaro .

ENV DB_URL="add-your-own-db-url"

EXPOSE 8008

CMD ["/root/jaro"]
