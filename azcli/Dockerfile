FROM alpine:3.9

ADD azgo / 


RUN \
  apk update && \
  apk add bash py-pip && \
  apk add --virtual=build gcc libffi-dev musl-dev openssl-dev python-dev make && \
  pip --no-cache-dir install -U pip && \
  pip install azure-cli && \
  apk del --purge build
  
CMD ["/azgo"]