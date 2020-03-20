FROM golang:alpine AS build

WORKDIR $GOPATH/src/github.com/felipefrizzo/brazilian-zipcode-api

COPY . $GOPATH/src/github.com/felipefrizzo/brazilian-zipcode-api
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /usr/bin/brazilian-zipcode-api

FROM alpine
LABEL MAINTAINER="Felipe Frizzo felipefrizzo@gmail.com"

RUN apk add --update --no-cache ca-certificates

COPY --from=build /usr/bin/brazilian-zipcode-api /usr/bin/brazilian-zipcode-api
EXPOSE 8000
ENTRYPOINT [ "/usr/bin/brazilian-zipcode-api" ]