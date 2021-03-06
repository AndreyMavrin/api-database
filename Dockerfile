FROM golang:1.15.2-buster AS build

RUN mkdir /app

COPY . /app

WORKDIR /app

RUN go build -o api-server main.go

FROM ubuntu:20.04

RUN apt-get -y update && apt-get install -y tzdata

# install Postgres
ENV PGVER 12
RUN apt-get -y update && apt-get install -y postgresql-$PGVER

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
RUN echo "listen_addresses='*'\nsynchronous_commit = off\nfsync = off\n" >> /etc/postgresql/$PGVER/main/postgresql.conf

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