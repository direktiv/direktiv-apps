FROM golang:1.15-buster as build

WORKDIR /go/src/app
ADD . /go/src/app
RUN go get -d -v
RUN CGO_ENABLED=0 go build -o /isolate -ldflags="-s -w" main.go

FROM ubuntu
RUN apt update -y
RUN apt install -y ansible
COPY --from=build /isolate /
CMD ["/isolate"]
