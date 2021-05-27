# FROM alpine:latest as build-terraform
# RUN wget https://releases.hashicorp.com/terraform/0.15.3/terraform_0.15.3_linux_amd64.zip
# RUN unzip terraform_0.15.3_linux_amd64.zip
# RUN ./terraform -help

# FROM alpine:latest as get-certs
# RUN apk --update add ca-certificates
# RUN mkdir /test

FROM golang:1.16-buster as build
WORKDIR /go/src/app
ADD . /go/src/app
RUN go get -d -v
RUN CGO_ENABLED=0 go build -o /terraform-wrap -ldflags="-s -w" main.go

FROM alpine:latest
# terraform
RUN wget https://releases.hashicorp.com/terraform/0.15.3/terraform_0.15.3_linux_amd64.zip
RUN unzip terraform_0.15.3_linux_amd64.zip
RUN ./terraform -help

# certs
RUN apk --update add ca-certificates

#COPY --from=build-terraform /terraform /
#COPY --from=get-certs /test /tmp
#COPY --from=get-certs /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=build /terraform-wrap /
CMD ["/terraform-wrap"]
