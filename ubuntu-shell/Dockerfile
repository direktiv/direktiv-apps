FROM golang:1.15-buster as build

WORKDIR /go/src/app
ADD . /go/src/app
RUN go get -d -u -v
RUN CGO_ENABLED=0 go build -o /shellrunner -ldflags="-s -w" main.go

FROM ubuntu
RUN apt-get update -y
RUN apt-get install jq wget curl git perl build-essential  -y
COPY --from=build /shellrunner /
CMD ["/shellrunner"]
