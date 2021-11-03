# FROM golang as build

# WORKDIR /app
# COPY ./main.go ./
# COPY ./go.mod ./
# COPY ./go.sum ./
# RUN go get -u -v

FROM alpine:latest as certs
# COPY --from=build /buildgo /

RUN apk add --no-cache git make musl-dev go

ENV GOROOT /usr/lib/go
ENV GOPATH /go
ENV PATH /go/bin:$PATH

WORKDIR /app
COPY ./main.go ./
COPY ./go.mod ./
COPY ./go.sum ./
RUN mkdir -p ${GOPATH}/src ${GOPATH}/bin
RUN CGO_ENABLED=0 go build -o /buildgo -ldflags="-s -w" main.go

CMD [ "/buildgo" ]