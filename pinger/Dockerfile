FROM golang:1.23-alpine

WORKDIR /app

COPY ./common common
COPY ./pinger/go.mod ./pinger/go.sum pinger/

WORKDIR /app/pinger

RUN go mod download

COPY ./pinger .


RUN go build -o main ./cmd/pinger

CMD ["./main"]