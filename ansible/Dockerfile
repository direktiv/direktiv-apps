FROM golang:1.15-buster as build

WORKDIR /go/src/app
ADD . /go/src/app
RUN go get -d -v
RUN CGO_ENABLED=0 go build -o /isolate -ldflags="-s -w" main.go

FROM ansible/ubuntu14.04-ansible
RUN apt update -y
RUN apt install -y python
COPY --from=build /isolate /
CMD ["/isolate"]
