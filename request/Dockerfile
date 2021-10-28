FROM golang as build

WORKDIR /app
COPY ./main.go ./
COPY ./go.mod ./
COPY ./go.sum ./
RUN go get -u -v
RUN CGO_ENABLED=0 go build -o /request -ldflags="-s -w" main.go

FROM alpine:latest as certs
RUN apk --update add ca-certificates

FROM scratch
COPY --from=build /request /
COPY --from=certs /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
CMD ["/request"]
