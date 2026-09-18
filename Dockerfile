
FROM golang:1.21-alpine AS builder
WORKDIR /app
# ONLY copy go.mod. Do NOT copy go.sum because it is deleted and causing build failures.
COPY go.mod ./
# Force resolve dependencies
RUN go mod tidy
# Copy the rest of the source code
COPY . .
# Build the binary
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/api/main.go

FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/main .
EXPOSE 8080
CMD ["./main"]
