ARG BASE
FROM siose-innova/gdal:3.0.2 AS etl
COPY ./src /usr/src/pg_runoff
ARG CS="CREATE_SCHEMA=OFF"
ARG CT="CREATE_TABLE=OFF"
ARG DT="DROP_TABLE=OFF"
ARG I1="SELECT COD_CORINE AS clc_key, CORINE_TXT clc_desc FROM data GROUP BY COD_CORINE, CORINE_TXT"
ARG I2="SELECT COD_MOPU AS mopu_key, MOPU_TXT mopu_desc FROM data GROUP BY COD_MOPU, MOPU_TXT"
ARG I3="SELECT COD_SIOSE AS siose_key, COD_CORINE AS clc_key, COD_MOPU AS mopu_key, MANNING AS n FROM data"
RUN set -ex \
    \
    && cd /usr/src/pg_runoff \
    && ogr2ogr -f "PGDUMP" --config PG_USE_COPY YES -dialect SQLITE \
               -sql "$I1" -nln clc_keys -lco SCHEMA=@extschema@ -lco $CS -lco $CT -lco $DT \
               data1.sql data.csv \
    && sed -i 's/"\@extschema\@"/\@extschema\@/' data1.sql \
    && ogr2ogr -f "PGDUMP" --config PG_USE_COPY YES -dialect SQLITE \
               -sql "$I2" -nln mopu_keys -lco SCHEMA=@extschema@ -lco $CS -lco $CT -lco $DT \
               data2.sql data.csv \
    && sed -i 's/"\@extschema\@"/\@extschema\@/' data2.sql \
    && ogr2ogr -f "PGDUMP" --config PG_USE_COPY YES -dialect SQLITE \
               -sql "$I3" -nln siose_keys -lco SCHEMA=@extschema@ -lco $CS -lco $CT -lco $DT \
               data3.sql data.csv \
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
