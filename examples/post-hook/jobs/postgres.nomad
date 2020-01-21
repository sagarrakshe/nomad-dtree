job "postgres" {
  datacenters = ["dc1"]
  type        = "service"

  group "postgres" {
    count = 1

    task "postgres" {
      driver = "docker"

      env {
        POSTGRES_DB       = "master"
        POSTGRES_USER     = "admin"
        POSTGRES_PASSWORD = "d8851a7b4720"
        PGDATA            = "/var/lib/postgresql/data/pgdata"
      }

      template {
        data = <<EOH
local   all             all                                     trust
host    all             all             127.0.0.1/32            trust
host    all             all             ::1/128                 trust
local   replication     all                                     trust
host    replication     all             127.0.0.1/32            trust
host    replication     all             ::1/128                 trust
host    all             all             all                     md5
host    replication     all             0.0.0.0/0               md5
        EOH

        destination = "/etc/postgresconf/pg_hba.conf"
        perms       = "444"
      }

      template {
        data = <<EOH
listen_addresses = '*'
max_connections = 100
shared_buffers = 128MB
dynamic_shared_memory_type = posix
log_timezone = 'UTC'
datestyle = 'iso, mdy'
timezone = 'UTC'
lc_messages = 'en_US.utf8'
lc_monetary = 'en_US.utf8'
lc_numeric = 'en_US.utf8'
lc_time = 'en_US.utf8'
default_text_search_config = 'pg_catalog.english'
wal_level = hot_standby
max_wal_senders = 8
wal_keep_segments = 8
hot_standby = on
        EOH

        destination = "/etc/postgresconf/postgresql.conf"
        perms       = "444"
      }

      
template {
    data = <<EOH
nameserver {{ env "attr.unique.network.ip-address" }}
EOH

    destination = "etc/resolv.conf"
}

      config {
        image = "postgres:12"

        port_map {
          db = 5432
        }

        volumes = [
          "etc/postgresconf/pg_hba.conf:/var/lib/postgresql/data/pg_hba.conf",
          "etc/postgresconf/postgresql.conf:/var/lib/postgresql/data/postgresql.conf",
          "etc/resolv.conf:/etc/resolv.conf"
        ]
      }

      resources {
        cpu    = 50
        memory = 256

        network {
          port "db" {
            static = "5432"
          }
        }
      }
    }
  }
}
