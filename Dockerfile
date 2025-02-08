# Stage 1: Build the Go application
FROM golang:1.21 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o movie-review-app .

# Stage 2: Create a minimal production image
FROM alpine:latest

WORKDIR /root/

COPY --from=builder /app/movie-review-app .

EXPOSE 8080

CMD ["./movie-review-app"]