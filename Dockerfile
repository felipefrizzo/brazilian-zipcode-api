FROM golang:alpine AS builder

WORKDIR $GOPATH/src/github.com/felipefrizzo/brazilian-zipcode-api
RUN apk add --update --no-cache ca-certificates

COPY . $GOPATH/src/github.com/felipefrizzo/brazilian-zipcode-api
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /usr/bin/brazilian-zipcode-api

FROM scratch
LABEL MAINTAINER="Felipe Frizzo felipefrizzo@gmail.com"

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=builder /usr/bin/brazilian-zipcode-api /usr/bin/brazilian-zipcode-api
EXPOSE 8000
ENTRYPOINT [ "/usr/bin/brazilian-zipcode-api" ]