FROM golang:1.15-buster as build

WORKDIR /go/src/app
ADD . /go/src/app
RUN go get -d -v 
RUN CGO_ENABLED=0 go build -o /smtp-bare -ldflags="-s -w" main.go

FROM alpine:latest as certs
RUN apk --update add ca-certificates

FROM scratch
COPY --from=build /smtp-bare /
COPY --from=certs /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
CMD ["/smtp-bare"]