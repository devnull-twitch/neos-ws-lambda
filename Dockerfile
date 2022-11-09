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
ENV ASSET_DIR assets
COPY assets assets
ENV HTML_TEMPLATE_DIR templates
COPY templates templates
ENV LAMBDA_TEMPLATE_DIR session_templates
RUN mkdir session_templates

ENTRYPOINT ["/lambda-server"]