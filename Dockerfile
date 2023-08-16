FROM golang:1.19.7-alpine3.17 as builder
RUN apk add --no-cache gcc git make musl-dev

WORKDIR /app
COPY . .
RUN go build -o bin/tx-generator ./tx-generator/cmd/main.go
RUN go build -o bin/malicious-validator ./validators/cmd/main.go

FROM alpine:3.17 as runner

RUN addgroup user && \
    adduser -G user -s /bin/sh -h /home/user -D user

USER user
WORKDIR /home/user/

FROM runner as malicious-validator
COPY --from=builder /app/bin/malicious-validator /usr/local/bin
ENTRYPOINT ["malicious-validator"]

FROM runner as tx-generator
COPY --from=builder /app/bin/tx-generator /usr/local/bin
ENTRYPOINT ["tx-generator"]
