FROM golang:1.17-alpine AS builder

RUN apk add --no-cache ca-certificates git

WORKDIR /app

COPY . .

COPY ./netrc /root/.netrc
RUN chmod 600 /root/.netrc

RUN go mod download && \
    go build -o lambda-server cmd/server/main.go

FROM alpine:3 AS runner

COPY --from=builder /app/lambda-server /lambda-server

ENTRYPOINT ["/lambda-server"]