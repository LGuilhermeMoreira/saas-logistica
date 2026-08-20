FROM golang:1.26-alpine AS builder
WORKDIR /app

COPY proto/ ./proto/
COPY auth/ ./auth/

WORKDIR /app/auth
RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux go build -o /app/server ./cmd

FROM alpine:latest
RUN apk --no-cache add ca-certificates tzdata
WORKDIR /app

COPY --from=builder /app/server .
RUN chmod +x /app/server

EXPOSE 50051

CMD ["/app/server"]
