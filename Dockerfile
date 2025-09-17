FROM golang:alpine3.22 AS builder

WORKDIR /usr/src/app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o hangout-content-delivery-api .

FROM alpine:3

RUN addgroup -S appgroup && adduser -S appuser -G appgroup

WORKDIR /home/appuser/app

COPY --chown=appuser:appgroup --from=builder /usr/src/app/resources/application.yaml ./resources/
COPY --chown=appuser:appgroup --from=builder /usr/src/app/hangout-content-delivery-api ./

RUN mkdir -p /mnt/certs && chown -R appuser:appgroup /mnt/certs

USER appuser

CMD ["./hangout-content-delivery-api"]
