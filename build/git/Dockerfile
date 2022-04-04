FROM golang:1.17-buster as builder
COPY go.mod src/
COPY go.sum src/
RUN cd src/ && go mod download
COPY cmd src/cmd/
RUN cd src && GO_ENABLED=0 go build -o /application -ldflags=" -s -w" cmd/git/*.go

FROM ubuntu:21.04
RUN apt-get update
RUN apt install openssh-client -y
RUN apt install git -y
COPY --from=builder /application /application
CMD ["/application"]
