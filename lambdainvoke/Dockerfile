FROM scratch

COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
ADD lambda /

CMD ["/lambda"]
