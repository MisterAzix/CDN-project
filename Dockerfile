FROM golang:1.22.3

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod tidy

COPY . .

RUN go build -o server

EXPOSE 8080

CMD ["/app/server"]