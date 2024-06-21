FROM golang:1.22-alpine3.20 as builder
WORKDIR /golang-mqtt-producer 
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" .
FROM busybox
WORKDIR /golang-mqtt-producer
COPY --from=builder /golang-mqtt-producer  /usr/bin/
ENTRYPOINT ["golang-mqtt-producer"]