FROM golang:1.26 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

ENV CGO_ENABLED=0

RUN go build -o smtp-webhook-server .

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/smtp-webhook-server /usr/local/bin/smtp-webhook-server

EXPOSE 2525

ENTRYPOINT ["smtp-webhook-server"]

CMD [
"--listen=:2525",
"--timeout.read=60",
"--timeout.write=60"
]
