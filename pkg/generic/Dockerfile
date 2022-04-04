FROM golang:1.17-buster as builder
COPY go.mod src/
COPY go.sum src/
RUN cd src/ && go mod download
COPY main.go src/cmd/
RUN cd src && GO_ENABLED=0 go build -o /direktiv-generic -tags osusergo,netgo -ldflags=" -s -w" cmd/*.go

RUN mkdir /out
CMD ["cp", "/direktiv-generic", "/out"]
