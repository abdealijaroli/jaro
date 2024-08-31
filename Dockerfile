FROM golang:1.22.3 AS builder

WORKDIR /app

LABEL maintainer="Abdeali Jaroli <abdeali@jaro.li>"

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o jaro .

FROM alpine:latest

WORKDIR /root/

COPY --from=builder /app/jaro .

EXPOSE 8008

CMD ["./jaro"]