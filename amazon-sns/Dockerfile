FROM scratch

ADD ./amazonsns/amazon-sns / 
COPY /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt

CMD ["amazon-sns"]