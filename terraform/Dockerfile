FROM golang as build


WORKDIR /go/src/app
ADD . /go/src/app
RUN go get -d -u -v
RUN CGO_ENABLED=0 go build -o /terraform-wrap -ldflags="-s -w" main.go

FROM alpine:latest
# terraform
RUN wget https://releases.hashicorp.com/terraform/0.15.3/terraform_0.15.3_linux_amd64.zip
RUN unzip terraform_0.15.3_linux_amd64.zip
RUN ./terraform -help

# certs
RUN apk --update add ca-certificates

#git
RUN apk add --no-cache bash git openssh

COPY --from=build /terraform-wrap /
CMD ["/terraform-wrap"]
