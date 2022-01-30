FROM golang:1.16-alpine AS builder
WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -mod vendor -ldflags="-w -s" -o /bin/log-analyzer ./cmd/...

FROM alpine:latest AS log-analyzer
RUN apk add curl && apk add jq
COPY --from=builder /bin/log-analyzer /usr/local/bin/log-analyzer
EXPOSE 4200
CMD ["/usr/local/bin/log-analyzer"]