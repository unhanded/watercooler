FROM golang:1.22.2-alpine3.19 AS builder

WORKDIR /app

COPY . .

RUN go mod download

RUN go build ./app/watercooler-serve

FROM alpine:3.19 AS RUNNER

WORKDIR /app

COPY --from=builder /app/watercooler-serve .

CMD ["./watercoolr-serve"]

EXPOSE 8000
