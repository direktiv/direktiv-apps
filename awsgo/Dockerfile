FROM golang:1.15-buster as build

WORKDIR /app
COPY ./main.go ./
COPY ./go.mod ./
COPY ./go.sum ./
RUN go get -u -v
RUN CGO_ENABLED=0 go build -o /awsgo -ldflags="-s -w" main.go

FROM alpine:3.9
RUN \
  apk update && \
  apk add --no-cache ca-certificates && \
  apk add bash python3 py3-pip && \
  apk add --virtual=build gcc libffi-dev musl-dev openssl-dev python3-dev make && \
  python3 --version && \
  python3 -m pip --no-cache-dir install -U pip && \
  python3 -m pip --no-cache-dir install awscli && \
  apk del --purge build

COPY --from=build /awsgo /

CMD ["/awsgo"]