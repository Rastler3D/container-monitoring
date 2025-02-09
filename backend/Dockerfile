FROM golang:1.23-alpine

WORKDIR /app

COPY ./common common
COPY ./backend/go.mod ./backend/go.sum backend/

WORKDIR /app/backend

RUN go mod download

COPY ./backend .

RUN go build -o main ./cmd/backend

CMD ["./main"]

