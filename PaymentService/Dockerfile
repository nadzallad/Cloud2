FROM golang:1.25

WORKDIR /app

COPY go.mod ./
RUN go mod tidy

COPY . .

WORKDIR /app/PaymentService
RUN go build -o app

CMD ["./app"]