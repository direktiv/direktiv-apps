FROM golang:1.16.13-buster as build

WORKDIR /app
COPY ./main.go ./
COPY ./go.mod ./
COPY ./go.sum ./

RUN go get -u -v
RUN CGO_ENABLED=0 go build -o /network -ldflags="-s -w" main.go

FROM alpine:latest as certs
RUN apk --update add ca-certificates

FROM scratch

COPY --from=build /network /
COPY --from=certs /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
CMD ["/network"]