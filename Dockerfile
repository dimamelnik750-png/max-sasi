FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o app .

FROM alpine:latest

WORKDIR /

RUN apk --no-cache add ca-certificates

COPY --from=builder /app/app .
COPY --from=builder /app/.env .

EXPOSE 8080

CMD ["./app"]