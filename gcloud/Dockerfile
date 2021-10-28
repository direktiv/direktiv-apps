FROM golang:1.15-buster as build

WORKDIR /app
COPY ./main.go ./
COPY ./go.mod ./
COPY ./go.sum ./
RUN go get -u -v
RUN CGO_ENABLED=0 go build -o /gcli -ldflags="-s -w" main.go

FROM alpine:3.9
 
RUN apk add --update python curl which bash 
RUN curl -sSL https://sdk.cloud.google.com | bash 
ENV PATH $PATH:/root/google-cloud-sdk/bin

COPY --from=build /gcli /

CMD ["/gcli"]