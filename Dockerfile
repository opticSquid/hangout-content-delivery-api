FROM golang:alpine3.22 AS builder

WORKDIR /usr/src/app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o hangout-content-delivery-api .

FROM alpine:3

RUN addgroup -S appgroup && adduser -S appuser -G appgroup

WORKDIR /app

COPY --from=builder /usr/src/app/hangout-content-delivery-api .

RUN chown appuser:appgroup /app/hangout-content-delivery-api

USER appuser

CMD ["./hangout-content-delivery-api"]