job "redis" {
  datacenters = ["dc1"]
  type        = "service"

  group "redis" {
    count = 1

    task "redis" {
      driver = "docker"

      template {
        data = <<EOH
dir /data
appendonly yes
appendfilename "appendonly.aof"
appendfsync everysec
EOH

        destination = "etc/redis.conf"
      }

template {
    data = <<EOH
nameserver {{ env "attr.unique.network.ip-address" }}
EOH

    destination = "etc/resolv.conf"
}

      config {
        image = "redis:5"

        port_map {
          listen = 6379
        }

        volumes = [
          "etc/redis.conf:/etc/redis/redis.conf",
          "etc/resolv.conf:/etc/resolv.conf",
        ]

        entrypoint = ["redis-server", "/etc/redis/redis.conf"]
      }

      resources {
        cpu    = 500
        memory = 4096

        network {

          port "listen" {
            static = "6379"
          }
        }
      }
    }
  }
}
