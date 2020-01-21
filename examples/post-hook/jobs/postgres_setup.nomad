job "postgres_setup" {
  datacenters = ["dc1"]
  type        = "batch"

  group "db_config" {
    count = 1

    task "initdb" {
      driver = "docker"

      env {
        CONN_STRING = "postgresql://admin:d8851a7b4720@${attr.unique.network.ip-address}:5432/master"
        DATABASE    = "testdb"
        USER        = "testuser"
        PASSWORD    = "testpass@123"
      }

template {
    data = <<EOH
nameserver {{ env "attr.unique.network.ip-address" }}
EOH

    destination = "etc/resolv.conf"
}

      template {
        data = <<EOH
#! /bin/bash

(
cat <<-EOF
CREATE ROLE ${USER} LOGIN PASSWORD '${PASSWORD}';
CREATE DATABASE ${DATABASE};
\c ${DATABASE};
CREATE EXTENSION IF NOT EXISTS citext;
GRANT ALL PRIVILEGES ON DATABASE ${DATABASE} TO ${USER};

EOF
) > /etc/initdb.sql

psql ${CONN_STRING} < /etc/initdb.sql
EOH

        destination = "/etc/scripts/initdb.sh"
        perms       = "111"
      }

      config {
        image   = "postgres:12"
        command = "/etc/initdb.sh"

        volumes = [
          "etc/scripts/initdb.sh:/etc/initdb.sh",
          "etc/resolv.conf:/etc/resolv.conf"
        ]
      }

      resources {
        cpu = 50
        memory = 256
      }
    }
  }
}
