FROM golang:1.19 as build

WORKDIR /app
COPY ./main.go ./
COPY ./go.mod ./
COPY ./go.sum ./

RUN go get -u -v
RUN CGO_ENABLED=0 go build -o /network -ldflags="-s -w" main.go

FROM alpine:latest as certs
RUN apk --update add ca-certificates

FROM alpine

RUN apk update && apk add bind-tools

COPY --from=build /network /
COPY --from=certs /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
CMD ["/network"]