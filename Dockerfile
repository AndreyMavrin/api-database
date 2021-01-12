FROM golang:1.15.2-buster AS build

RUN mkdir /app

COPY . /app

WORKDIR /app

RUN go build -o api-server main.go

FROM ubuntu:20.04 AS release

# Make the "en_US.UTF-8" locale so postgres will be utf-8 enabled by default
RUN apt-get -y update && apt-get install -y locales gnupg2 tzdata
RUN locale-gen en_US.UTF-8
RUN update-locale LANG=en_US.UTF-8

ENV TZ=Russia/Moscow
RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone

# install Postgres
ENV PGVER 12
RUN apt-get update -y && apt-get install -y postgresql postgresql-contrib

# Run the rest of the commands as the ``postgres`` user created by the ``postgres-$PGVER`` package when it was ``apt installed``
USER postgres

# Create a PostgreSQL role named ``amavrin`` with ``root`` as the password and
# then create a database `forums` owned by the ``amavrin`` role.
RUN /etc/init.d/postgresql start &&\
    psql --command "CREATE USER amavrin WITH SUPERUSER PASSWORD 'root';" &&\
    createdb -E UTF8 forums &&\
    /etc/init.d/postgresql stop

RUN echo "host all  all    0.0.0.0/0  md5" >> /etc/postgresql/$PGVER/main/pg_hba.conf

# Expose the PostgreSQL port
RUN echo "listen_addresses='*'\nsynchronous_commit = off\nfsync = off\nshared_buffers = 256MB\neffective_cache_size = 1536MB\n" >> /etc/postgresql/$PGVER/main/postgresql.conf
RUN echo "wal_buffers = 16MB\nwal_writer_delay = 50ms\nrandom_page_cost = 1.1\nmax_connections = 100\nwork_mem = 20971kB\nmaintenance_work_mem = 512MB\ncpu_tuple_cost = 0.0030\ncpu_index_tuple_cost = 0.0010\ncpu_operator_cost = 0.0005" >> /etc/postgresql/$PGVER/main/postgresql.conf
RUN echo "full_page_writes = off" >> /etc/postgresql/$PGVER/main/postgresql.conf
RUN echo "log_statement = none" >> /etc/postgresql/$PGVER/main/postgresql.conf
RUN echo "log_duration = off " >> /etc/postgresql/$PGVER/main/postgresql.conf
RUN echo "log_lock_waits = on" >> /etc/postgresql/$PGVER/main/postgresql.conf
RUN echo "log_min_duration_statement = 5000" >> /etc/postgresql/$PGVER/main/postgresql.conf
RUN echo "log_filename = 'query.log'" >> /etc/postgresql/$PGVER/main/postgresql.conf
RUN echo "log_directory = '/var/log/postgresql'" >> /etc/postgresql/$PGVER/main/postgresql.conf
RUN echo "log_destination = 'csvlog'" >> /etc/postgresql/$PGVER/main/postgresql.conf
RUN echo "logging_collector = on" >> /etc/postgresql/$PGVER/main/postgresql.conf
RUN echo "log_temp_files = '-1'" >> /etc/postgresql/$PGVER/main/postgresql.conf

EXPOSE 5432

# Add VOLUMEs to allow backup of config, logs and databases
VOLUME ["/etc/postgresql", "/var/log/postgresql", "/var/lib/postgresql"]

# Back to the root user
USER root

# Собранный сервер
COPY --from=build /app /app

EXPOSE 5000
ENV PGPASSWORD root
CMD service postgresql start && psql -h localhost -d forums -U amavrin -p 5432 -a -q -f /app/storage/migrations/up.sql && /app/api-server 