FROM alpine:3.6 
RUN apk add --update python curl which bash 
RUN curl -sSL https://sdk.cloud.google.com | bash 
ENV PATH $PATH:/root/google-cloud-sdk/bin
