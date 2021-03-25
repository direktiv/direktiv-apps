FROM scratch

ADD ./amazonsns/amazon-sns / 
RUN mkdir -p /etc/ssl/certs 
ADD ./amazonsns/ca-certificates.crt /etc/ssl/certs/

CMD ["amazon-sns"]