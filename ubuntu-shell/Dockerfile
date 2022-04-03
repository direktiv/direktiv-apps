FROM golang:1.17-buster as build

WORKDIR /app
COPY ./go.mod ./
COPY ./go.sum ./
RUN go mod download

COPY ./main.go ./
RUN CGO_ENABLED=0 go build -o /shellrunner -ldflags="-s -w" main.go

FROM ubuntu:21.10

ARG DEBIAN_FRONTEND=noninteractive
RUN apt-get update -y
RUN apt-get install sed jq wget curl git perl \
        build-essential openssh-server openssh-client \
        golang-1.17 python3 python3-pip make -y
COPY --from=build /shellrunner /
CMD ["/shellrunner"]
