FROM golang:1.25

WORKDIR /app

COPY . .

RUN go mod tidy

RUN go build -o app .

EXPOSE 8085

CMD ["./app"]
