FROM alpine:3.10

RUN set -ex \
    \
    && apk add --no-cache \
        bash \
        bash-completion \
        postgresql-client \
        tzdata \
    && cp /usr/share/zoneinfo/Europe/Madrid /etc/localtime

ENTRYPOINT [ "psql" ]
