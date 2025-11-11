# Build stage
FROM golang:1.24.1-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git ca-certificates

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/main.go

# Runtime stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates wget

RUN addgroup -g 1001 -S appgroup \
  && adduser -u 1001 -S appuser -G appgroup

WORKDIR /app

COPY --from=builder /app/main .

# ✅ make sure it's executable and owned by our non-root user
RUN chmod +x ./main && chown -R appuser:appgroup /app

USER appuser

EXPOSE 9000
CMD ["./main"]



