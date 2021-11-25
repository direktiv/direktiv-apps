FROM golang:1.17.3-buster as build

WORKDIR /app
COPY ./main.go ./
RUN go mod init github.com/direktiv-apps/jmeter
COPY ./go.mod ./
COPY ./go.sum ./
RUN go mod tidy
RUN go get -u -v
RUN CGO_ENABLED=0 go build -o /jmeter -ldflags="-s -w" main.go


FROM ubuntu:21.10

RUN apt-get update && apt-get install openjdk-17-jdk jmeter -y
COPY --from=build /jmeter /


CMD ["/jmeter"]
