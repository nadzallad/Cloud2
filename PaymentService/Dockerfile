FROM golang:1.25

WORKDIR /app

# copy folder service ke dalam container
COPY PaymentService/ PaymentService/

WORKDIR /app/PaymentService

RUN go mod tidy
RUN go build -o app .

CMD ["./app"]