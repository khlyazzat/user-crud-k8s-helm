FROM golang:1.24.3-alpine AS builder
WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o userservice

FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/userservice .
ENTRYPOINT ["./userservice"]
