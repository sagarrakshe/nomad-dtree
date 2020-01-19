job "nginx" {
  datacenters = ["dc1"]
  type        = "service"

  group "nginx" {
    count = 1

    task "nginx" {
      driver = "docker"

      template {
        data = <<EOF
user nginx;
worker_processes  5;  ## Default: 1

daemon off;

events {
  worker_connections  4096;  ## Default: 1024
}

http {
  # Some basic config.
  server_tokens off;
  sendfile      on;
  tcp_nopush    on;
  tcp_nodelay   on;

resolver {{ env "attr.unique.network.ip-address" }} valid=5s;

server {
    listen 28888;

    server_name localhost.com;

    access_log /dev/stdout;
    error_log /dev/stderr;

    location / {
      proxy_pass http://{{ env "attr.unique.network.ip-address" }}:8888;
      proxy_send_timeout 900s;
      proxy_set_header Host $http_host;
      proxy_set_header X-Real-IP $remote_addr;
      proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
      proxy_set_header X-Forwarded-Proto $scheme;
      add_header         Access-Control-Allow-Headers 'Authorization, Content-Type' always;
    }

  }
}
EOF
        destination = "/etc/nginx_root.conf"
        perms       = "755"
      }

template {
    data = <<EOH
nameserver {{ env "attr.unique.network.ip-address" }}
EOH

    destination = "etc/resolv.conf"
}

    config {
      image = "nginx:1.15.11"
      network_mode = "host"

      command = "/bin/bash"
      args    = ["-c", "nginx -c /etc/nginx/nginx.conf"]


      volumes = [
          "etc/nginx_root.conf:/etc/nginx/nginx.conf",
          "certs/:/etc/nginx/certs/public/",
          "etc/resolv.conf:/etc/resolv.conf",
        ]

      port_map {
        http = 28888
      }
    }

    resources {
        cpu    = 500
        memory = 1024

        network {
          port "http" {
            static = 28888
          }
        }
      }
    }
  }
}
