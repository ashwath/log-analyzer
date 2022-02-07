##
# Downloads Golang and builds log-analyzer
##
FROM golang:1.16-alpine AS builder
WORKDIR /app
COPY . .
RUN apk add build-base && apk add jq
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -mod vendor -ldflags="-w -s" -o /bin/log-analyzer ./cmd/...
RUN go test --mod=vendor -v -coverprofile=coverage.out -tags=unit -json ./...

##
# Starts conatiner log-analyzer using output of @builder
##
FROM alpine:latest AS log-analyzer

## Need jq to parse the json output
RUN apk add curl && apk add jq

COPY --from=builder /bin/log-analyzer /usr/local/bin/log-analyzer

## Copies test files to /var/log of the conatiner
COPY test-files /var/log/.

EXPOSE 4200
CMD ["/usr/local/bin/log-analyzer"]