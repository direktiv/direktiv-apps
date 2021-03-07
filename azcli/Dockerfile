FROM alpine

ADD azgo / 

RUN \
  apk update && \
  apk add bash py-pip && \
  apk add --no-cache ca-certificates && \
  apk add --virtual=build gcc libffi-dev musl-dev openssl-dev python-dev make && \
  pip --no-cache-dir install -U pip && \
  pip --no-cache-dir install azure-cli && \
  apk del --purge build

CMD ["/azgo"]