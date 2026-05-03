FROM golang:1.25

WORKDIR /app

# copy folder tracking service
COPY TrackingService/ TrackingService/

WORKDIR /app/TrackingService

RUN go mod tidy
RUN go build -o app .

CMD ["./app"]