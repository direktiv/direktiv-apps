FROM golang:1.17-buster as build

WORKDIR /app
COPY ./main.go ./
COPY ./go.mod ./
COPY ./go.sum ./
RUN go mod download
RUN CGO_ENABLED=0 go build -o /slackmsg -ldflags="-s -w" main.go

# FROM alpine:latest as certs
# RUN apk --update add ca-certificates

FROM gcr.io/distroless/static
USER nonroot:nonroot

COPY --from=build /slackmsg /
# COPY --from=certs /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
CMD ["/slackmsg"]