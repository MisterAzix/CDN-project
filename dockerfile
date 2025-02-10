FROM golang:1.22.2

WORKDIR /app

COPY go.mod ./
COPY go.mod go.sum ./

COPY app/ .

RUN go build -o server

EXPOSE 8080

CMD ["/app/server"]
