FROM alpine

ADD azgo / 

RUN apk update  
RUN apk add bash py-pip 
RUN apk add --no-cache ca-certificates
RUN apk add --virtual=build gcc libffi-dev musl-dev openssl-dev  make
RUN pip --no-cache-dir install -U pip
RUN pip --no-cache-dir install azure-cli
RUN apk del --purge build

CMD ["/azgo"]