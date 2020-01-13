FROM siose-innova/postgresql-client:11.6

COPY ext-create.sh /usr/local/bin
RUN chmod +x /usr/local/bin/ext-create.sh

ENTRYPOINT [ "/usr/local/bin/ext-create.sh" ]
