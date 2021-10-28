FROM golang:latest as build

WORKDIR /app
COPY ./main.go ./
COPY ./go.mod ./
COPY ./go.sum ./
RUN go get -u -v
RUN CGO_ENABLED=0 go build -o /vmware-powercli -ldflags="-s -w" main.go

FROM alpine:latest as certs
RUN apk --update add ca-certificates

FROM vmware/powerclicore:latest
WORKDIR /root
# RUN Add-PSSnapin "VMware.VimAutomation.Core" 
COPY --from=build /vmware-powercli /
COPY --from=certs /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
CMD ["/vmware-powercli"]