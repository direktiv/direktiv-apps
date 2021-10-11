FROM golang:1.15-buster as build

WORKDIR /go/src/app
ADD . /go/src/app
RUN go get -d -u -v
RUN CGO_ENABLED=0 go build -o /gitrun -ldflags="-s -w" main.go

FROM ubuntu:latest

RUN apt-get update
RUN apt install openssh-client -y
RUN apt install git -y

COPY --from=build /gitrun /
CMD ["/gitrun"]
