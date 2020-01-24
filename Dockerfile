ARG BASE
FROM siose-innova/gdal:3.0.2 AS etl
COPY ./src /usr/src/pg_runoff
RUN set -ex \
    \
    && cd /usr/src/pg_runoff \
    && ogr2ogr -f "PGDUMP" --config PG_USE_COPY YES \
               -nln siose2hydro -lco SCHEMA=@extschema@ \
               -lco CREATE_SCHEMA=OFF -lco CREATE_TABLE=OFF -lco DROP_TABLE=OFF \
               data1.sql hydro2siose.csv \
    && sed -i 's/"\@extschema\@"/\@extschema\@/' data1.sql \
    && ogr2ogr -f "PGDUMP" --config PG_USE_COPY YES \
               -nln manning -lco SCHEMA=@extschema@ \
               -lco CREATE_SCHEMA=OFF -lco CREATE_TABLE=OFF -lco DROP_TABLE=OFF \
               data2.sql manning.csv \
    && sed -i 's/"\@extschema\@"/\@extschema\@/' data2.sql \
    && ogr2ogr -f "PGDUMP" --config PG_USE_COPY YES \
               -nln p0 -lco SCHEMA=@extschema@ \
               -lco CREATE_SCHEMA=OFF -lco CREATE_TABLE=OFF -lco DROP_TABLE=OFF \
               data3.sql p0.csv \
    && sed -i 's/"\@extschema\@"/\@extschema\@/' data3.sql \
    && cat data1.sql data2.sql data3.sql > data.sql

FROM $BASE

ARG SCHEMA
ENV SIOSE_SCHEMA=$SCHEMA

COPY --from=etl /usr/src/pg_runoff /usr/src/pg_runoff/.
RUN set -ex \
    \
    && apk add --no-cache --virtual .build-deps \
        make \
    && cd /usr/src/pg_runoff \
    && make \
    && make install \
    && cd / \
    && rm -rf /usr/src/pg_runoff \
    && apk del .build-deps
