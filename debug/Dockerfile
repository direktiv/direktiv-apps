FROM golang:1.15-buster as build

WORKDIR /app
COPY ./main.go ./
COPY ./go.mod ./
COPY ./go.sum ./
RUN go get -u -v
RUN CGO_ENABLED=0 go build -o /debug -ldflags="-s -w" main.go


FROM scratch
COPY --from=build /debug /

CMD ["/debug"]