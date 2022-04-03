FROM golang:1.17-buster as build

WORKDIR /app

COPY ./go.mod ./
COPY ./go.sum ./
RUN go mod download

COPY ./main.go ./

RUN CGO_ENABLED=0 go build -o /application -ldflags="-s -w" main.go

FROM ubuntu:21.04

RUN apt update && apt install curl git -y
RUN curl -fsSL https://cli.github.com/packages/githubcli-archive-keyring.gpg | dd of=/usr/share/keyrings/githubcli-archive-keyring.gpg
RUN echo "deb [arch=$(dpkg --print-architecture) signed-by=/usr/share/keyrings/githubcli-archive-keyring.gpg] https://cli.github.com/packages stable main" | tee /etc/apt/sources.list.d/github-cli.list > /dev/null
RUN apt update
RUN apt install gh -y

EXPOSE 8080
COPY --from=build /application /
CMD ["/application"]

