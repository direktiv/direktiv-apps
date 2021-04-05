FROM golang:1.15-buster as build

WORKDIR /go/src/app
ADD . /go/src/app
RUN go get -d -v 
RUN CGO_ENABLED=0 go build -o /gcli -ldflags="-s -w" main.go

FROM alpine:3.9
 
RUN apk add --update python curl which bash 
RUN curl -sSL https://sdk.cloud.google.com | bash 
ENV PATH $PATH:/root/google-cloud-sdk/bin

COPY --from=build /gcli /

CMD ["/gcli"]