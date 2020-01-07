FROM golang:1.13-stretch AS builder

WORKDIR /usr/src/app

COPY go.mod .
RUN go mod download

COPY . .
RUN make build

FROM ubuntu:19.04
ENV DEBIAN_FRONTEND=noninteractive
ENV PGVER 11
ENV PORT 8080
ENV POSTGRES_HOST localhost
ENV POSTGRES_PORT 5432
ENV POSTGRES_DB postgres
ENV POSTGRES_USER d
ENV POSTGRES_PASSWORD 1234
EXPOSE $PORT

RUN apt-get update && apt-get install -y postgresql-$PGVER && apt-get install -y build-essential

USER postgres

RUN service postgresql start &&\
    psql --command "CREATE USER d WITH SUPERUSER PASSWORD '1234';" &&\
    service postgresql stop

VOLUME  ["/etc/postgresql", "/var/log/postgresql", "/var/lib/postgresql"]

COPY --from=builder /usr/src/app/ .

CMD service postgresql start && ./apiserver