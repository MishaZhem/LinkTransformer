FROM golang:1.24.3

WORKDIR /app

COPY . .

RUN go build -o analyticsservice ./cmd/main/main.go

CMD ["./analyticsservice"]
