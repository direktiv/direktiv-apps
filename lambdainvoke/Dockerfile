FROM alpine:3.12.4


RUN apk add --no-cache ca-certificates
ADD lambda /

CMD ["/lambda"]
