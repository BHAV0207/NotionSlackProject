FROM golang:1.24.4

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o ai-service ./cmd

EXPOSE 8005

CMD ["./ai-service"]
