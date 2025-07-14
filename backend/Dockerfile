FROM golang:1.24.2-alpine AS builder

ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o /app/server ./cmd/server


FROM alpine:latest

WORKDIR /app
COPY --from=builder /app/server .
EXPOSE 8080

CMD ["./server"]