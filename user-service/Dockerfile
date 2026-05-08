FROM golang:1.21-alpine AS builder

WORKDIR /app
COPY go.mod ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o user-service ./cmd/main.go

FROM alpine:latest

RUN addgroup -S app && adduser -S app -G app
WORKDIR /app
COPY --from=builder /app/user-service .
USER app

EXPOSE 9091
ENTRYPOINT ["./user-service"]
