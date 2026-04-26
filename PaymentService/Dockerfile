FROM golang:1.25

WORKDIR /app

COPY PaymentService/ .

RUN go mod tidy
RUN go build -o app .

CMD ["./app"]
