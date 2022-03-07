FROM golang:1.16.13-buster as build

WORKDIR /app
COPY ./main.go ./
COPY ./go.mod ./
COPY ./go.sum ./

RUN go get -u -v
RUN CGO_ENABLED=0 go build -o /run-python -ldflags="-s -w" main.go

FROM alpine:latest as certs
RUN apk --update add ca-certificates

FROM python:3.10.2-buster

RUN pip3 install requests cryptography

RUN apt-get update && apt-get install libssl-dev -y

COPY --from=build /run-python /
CMD ["/run-python"]
