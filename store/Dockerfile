FROM golang:1.15-buster as build

WORKDIR /go/src/app
ADD . /go/src/app
RUN go get -d -u -v 
RUN CGO_ENABLED=0 go build -o /store -ldflags="-s -w" googlepusher.go

FROM alpine:3.9 as certs
RUN apk --update add ca-certificates

FROM scratch
COPY --from=build /store /
COPY --from=certs /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
CMD ["/store"]

