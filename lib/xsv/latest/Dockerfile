FROM alpine:edge

RUN set -ex \
    \
    && apk add --no-cache \
        bash \
        bash-completion \
        tzdata \
        libgcc \
    \
    && cp /usr/share/zoneinfo/Europe/Madrid /etc/localtime \
    \
    && apk add --no-cache --virtual .fetch-deps \
        git \
    && apk add --no-cache --virtual .build-deps \
        --repository http://dl-cdn.alpinelinux.org/alpine/edge/community \
        cargo \
        rust \
    \
    && mkdir -p /usr/src \
    && cd /usr/src \
    && git clone git://github.com/BurntSushi/xsv \
    && cd xsv \
    && cargo build --release \
    && cp target/release/xsv /usr/bin/. \
    && cd / \
    && rm -rf /usr/src/xsv \
    \
    && apk del .fetch-deps .build-deps \
    && rm -rf /root/.cargo
