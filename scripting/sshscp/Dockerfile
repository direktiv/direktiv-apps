FROM golang:1.16.13-buster as build

WORKDIR /app
COPY ./main.go ./
COPY ./go.mod ./
COPY ./go.sum ./

RUN go mod download
RUN CGO_ENABLED=0 go build -o /application -ldflags="-s -w" main.go

FROM gcr.io/distroless/static
EXPOSE 8080
COPY --from=build /application /
CMD ["/application"]

