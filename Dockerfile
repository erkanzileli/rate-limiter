FROM golang:1.15-alpine3.12 AS builder

WORKDIR /build

RUN apk --no-cache add gcc g++ make git

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o app

FROM alpine:3.12

WORKDIR /app

RUN apk --no-cache add ca-certificates

COPY --from=builder /build/app ./rate-limiter

ENTRYPOINT /app/rate-limiter