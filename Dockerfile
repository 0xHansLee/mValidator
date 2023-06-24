FROM golang:1.19.7-alpine3.17 as builder
RUN apk add --no-cache gcc git make musl-dev
WORKDIR /app
COPY . .
RUN go build -o malicious-validator ./cmd/main.go 

FROM alpine:3.17 as malicious-validator
COPY --from=builder /app /usr/local/bin

ENTRYPOINT ["malicious-validator"]
