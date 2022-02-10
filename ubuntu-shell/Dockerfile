FROM golang:1.16-buster as build

WORKDIR /app
COPY ./main.go ./
COPY ./go.mod ./
COPY ./go.sum ./
RUN go get -u -v
RUN CGO_ENABLED=0 go build -o /shellrunner -ldflags="-s -w" main.go

FROM ubuntu

ARG DEBIAN_FRONTEND=noninteractive
RUN apt-get update -y
RUN apt-get install sed jq wget curl git perl build-essential openssh-server openssh-client -y
COPY --from=build /shellrunner /
CMD ["/shellrunner"]
