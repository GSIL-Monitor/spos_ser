FROM scratch
ADD ca-certificates.crt /etc/ssl/certs/
ADD spos-auth /
ADD basic_model.conf /
ADD basic_policy.csv /
CMD ["/spos-auth"]



