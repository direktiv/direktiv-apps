FROM golang:latest as build

WORKDIR /go/src/app
ADD . /go/src/app
RUN go get -d -u -v 
RUN CGO_ENABLED=0 go build -o /vmware-powercli -ldflags="-s -w" main.go

FROM alpine:latest as certs
RUN apk --update add ca-certificates

FROM vmware/powerclicore:latest
WORKDIR /root
# RUN Add-PSSnapin "VMware.VimAutomation.Core" 
COPY --from=build /vmware-powercli /
COPY --from=certs /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
CMD ["/vmware-powercli"]