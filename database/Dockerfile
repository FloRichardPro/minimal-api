FROM mysql as builder

RUN ["sed", "-i", "s/exec \"$@\"/echo \"not running $@\"/", "/usr/local/bin/docker-entrypoint.sh"]

ENV MYSQL_ROOT_PASSWORD=root

COPY --chown=mysql:mysql . /docker-entrypoint-initdb.d/
WORKDIR /docker-entrypoint-initdb.d/
RUN touch script.sql  && cat schema.sql >> script.sql && cat import.sql >> script.sql && cat user.sql >> script.sql && rm import.sql schema.sql user.sql

RUN ["/usr/local/bin/docker-entrypoint.sh", "mysqld", "--datadir", "/initialized-db"]

FROM mysql

COPY --from=builder /initialized-db /var/lib/mysql
