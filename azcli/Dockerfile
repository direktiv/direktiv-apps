FROM alpine:3.9

RUN \
  apk update && \
  apk add bash python3 py3-pip && \
  apk add --virtual=build gcc libffi-dev musl-dev openssl-dev python-dev make && \
  python3 --version && \
  python3 -m pip --no-cache-dir install -U pip && \
  python3 -m pip --no-cache-dir install azure-cli && \
  apk del --purge build

ADD azgo / 

CMD ["/azgo"]