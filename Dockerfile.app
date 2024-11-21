FROM golang:1.23.3-alpine3.20

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build cmd/main.go

CMD ["./main"]