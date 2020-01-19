job "api" {
  datacenters = ["dc1"]
  type        = "service"

  group "api" {
    count = 1

    task "api" {
      driver = "docker"

      env {
        REDIS_URL = "redis://${attr.unique.network.ip-address}:6379/0"
      }

template {
    data = <<EOH
nameserver {{ env "attr.unique.network.ip-address" }}
EOH

    destination = "etc/resolv.conf"
}

      config {
        image = "sagarrakshe/nomad-dtree-api:latest"

        port_map {
          api = 8888
        }

        volumes = [
          "etc/resolv.conf:/etc/resolv.conf",
        ]
      }

      resources {
        cpu    = 500
        memory = 4028

        network {
          port "api" {
            static = "8888"
          }
        }
      }
    }
  }
}
