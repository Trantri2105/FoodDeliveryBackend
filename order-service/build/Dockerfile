FROM golang:latest
WORKDIR /app
COPY . .
RUN go mod tidy
RUN CGO_ENABLED=0 GOOS=linux go build -o /order-srv ./cmd

CMD ["/order-srv"]