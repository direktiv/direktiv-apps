FROM golang:1.15-buster as build

WORKDIR /go/src/app
ADD . /go/src/app
RUN go get -d -u -v
RUN CGO_ENABLED=0 go build -o /watermark -ldflags="-s -w" main.go

FROM ubuntu:latest as ttf
RUN apt-get update && apt-get install -y fonts-liberation

FROM scratch
COPY --from=build /watermark /
COPY --from=ttf /usr/share/fonts/truetype/liberation/LiberationMono-Regular.ttf /usr/share/fonts/truetype/liberation/LiberationMono-Regular.ttf
CMD ["/watermark"]
