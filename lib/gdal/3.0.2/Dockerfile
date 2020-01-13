FROM alpine:3.11

RUN set -ex \
    \
    && apk add --no-cache \
        bash \
        bash-completion \
        gdal \
        tar \
        tzdata \
        wget \
    \
    && cp /usr/share/zoneinfo/Europe/Madrid /etc/localtime