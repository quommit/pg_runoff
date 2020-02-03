ARG BASE

FROM siose-innova/gdal:3.0.2 AS etl

ARG VERSION
ENV EXTENSION_VERSION=$VERSION

COPY ./src /usr/src/pg_runoff
RUN set -ex \
    \
    && cd /usr/src/pg_runoff \
    && ogr2ogr -f "PGDUMP" --config PG_USE_COPY NO \
               -sql "SELECT CAST(icn52 AS smallint), CAST(siose AS smallint) FROM hydro2siose" \
               -nln hydro2siose \
               -lco SCHEMA=@extschema@ \
               -lco CREATE_SCHEMA=OFF \
               -lco CREATE_TABLE=OFF \
               -lco DROP_TABLE=OFF \
               hydro2siose.sql hydro2siose.csv \
    && sed -i 's/"\@extschema\@"/\@extschema\@/' hydro2siose.sql \
    && sed -i -e '/^BEGIN/d' -e '/^COMMIT/d' hydro2siose.sql \
    && ogr2ogr -f "PGDUMP" --config PG_USE_COPY NO \
               -sql "SELECT CAST(siose AS smallint), CAST(n AS float) FROM manning" \
               -nln manning \
               -lco SCHEMA=@extschema@ \
               -lco CREATE_SCHEMA=OFF \
               -lco CREATE_TABLE=OFF \
               -lco DROP_TABLE=OFF \
               manning.sql manning.csv \
    && sed -i 's/"\@extschema\@"/\@extschema\@/' manning.sql \
    && sed -i -e '/^BEGIN/d' -e '/^COMMIT/d' manning.sql \
    && ogr2ogr -f "PGDUMP" --config PG_USE_COPY NO \
               -sql "SELECT CAST(icn52 AS smallint), CAST(slope_mod AS smallint), CAST(soil_mod AS smallint), CAST(p0 AS float) FROM p0" \
               -nln p0 \
               -lco SCHEMA=@extschema@ \
               -lco CREATE_SCHEMA=OFF \
               -lco CREATE_TABLE=OFF \
               -lco DROP_TABLE=OFF \
               p0.sql p0.csv \
    && sed -i 's/"\@extschema\@"/\@extschema\@/' p0.sql \
    && sed -i -e '/^BEGIN/d' -e '/^COMMIT/d' p0.sql \
    && cat hydro2siose.sql manning.sql p0.sql >> pg_runoff--${EXTENSION_VERSION}.sql.template


FROM $BASE

ARG DB
ENV SIOSE_DB=$DB
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
