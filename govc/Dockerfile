FROM golang as build

WORKDIR /app
COPY ./main.go ./
COPY ./go.mod ./
COPY ./go.sum ./
RUN go get -u -v
RUN CGO_ENABLED=0 go build -o /govc-run -ldflags="-s -w" main.go

FROM alpine as govc
# ENV path=/usr/local/bin
RUN apk --no-cache add curl
RUN curl -L -o - "https://github.com/vmware/govmomi/releases/latest/download/govc_$(uname -s)_$(uname -m).tar.gz" | tar -C /usr/local/bin -xvzf - govc

FROM scratch 

ENV PATH=/usr/local/bin
COPY --from=govc /usr/local/bin/govc /usr/local/bin/govc
COPY --from=govc /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=build /govc-run /

CMD [ "/govc-run" ]