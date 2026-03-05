FROM golang:1.23-alpine3.19 AS builder

WORKDIR /usr/src/app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /app/prosync ./cmd/main.go

FROM alpine:3.19

WORKDIR /app

COPY --from=builder /app/prosync .

COPY --from=builder /usr/src/app/web ./web

RUN mkdir -p /app/logs

EXPOSE 8000

CMD ["./prosync"]