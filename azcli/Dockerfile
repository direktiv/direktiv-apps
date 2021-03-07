FROM alpine:3.9

ADD azgo / 

RUN apk add --no-cache curl tar openssl sudo bash jq python3
RUN apk add --virtual=build gcc libffi-dev musl-dev openssl-dev make python3-dev
RUN pip3 install virtualenv &&\
    python3 -m virtualenv /azure-cli
ARG AZURE_CLI_VERSION=2.0.61
RUN /azure-cli/bin/python -m pip --no-cache-dir install azure-cli==${AZURE_CLI_VERSION}

CMD ["/azgo"]