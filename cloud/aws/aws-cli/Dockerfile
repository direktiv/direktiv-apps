FROM golang:1.17.8-buster as build

WORKDIR /app
COPY ./main.go ./
COPY ./go.mod ./
COPY ./go.sum ./

RUN go mod download
RUN CGO_ENABLED=0 go build -o /application -ldflags="-s -w" main.go

FROM amazon/aws-cli:2.4.23

COPY --from=build /application /
CMD ["/application"]
ENTRYPOINT ["/application"]