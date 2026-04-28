# pakai image Go
FROM golang:1.22

# set working directory
WORKDIR /app

# copy semua file
COPY . .

# build aplikasi
RUN go build -o order-service

# expose port
EXPOSE 8080

# run aplikasi
CMD ["./order-service"]