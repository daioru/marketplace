FROM golang:1.24.1 

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o api-gateway ./cmd/api-gateway

EXPOSE 8080

CMD ["/app/api-gateway"]