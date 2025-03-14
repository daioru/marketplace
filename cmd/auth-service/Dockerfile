FROM golang:1.24.1 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o auth-service ./cmd/auth-service

FROM alpine:latest

WORKDIR /root/

COPY --from=builder /auth-service .

RUN chmod +x /root/auth-service

EXPOSE 50051

CMD [ "/root/auth-service" ]