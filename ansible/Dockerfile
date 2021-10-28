FROM golang:1.15-buster as build

WORKDIR /app
COPY ./main.go ./
COPY ./go.mod ./
COPY ./go.sum ./
RUN go get -u -v
RUN CGO_ENABLED=0 go build -o /isolate -ldflags="-s -w" main.go

FROM ubuntu
RUN apt update -y
RUN apt install -y software-properties-common
RUN add-apt-repository --yes --update ppa:ansible/ansible
RUN apt install -y ansible
RUN ln -s /usr/bin/python3 /usr/bin/python

COPY --from=build /isolate /
CMD ["/isolate"]
