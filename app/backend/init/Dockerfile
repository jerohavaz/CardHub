FROM alpine:3.21

RUN apk add --no-cache openssl

COPY init/jwt-init.sh /usr/local/bin/jwt-init.sh

ENTRYPOINT ["/bin/sh", "/usr/local/bin/jwt-init.sh"]