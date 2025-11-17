FROM golang:1.22-alpine AS build

WORKDIR /app

# Allow the Go toolchain to auto-install the required version (>=1.24.1)
ENV GOTOOLCHAIN=auto

# Pre-fetch dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the source
COPY . .

# Build the chat service binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o chat-service ./cmd


FROM alpine:3.20

WORKDIR /app

# Copy binary
COPY --from=build /app/chat-service .

EXPOSE 8080

CMD ["./chat-service"]


