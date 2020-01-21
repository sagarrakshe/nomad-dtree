#!/bin/bash

CONSUL_ADDR="http://localhost:8500"
NOMAD_ADDR="http://localhost:4646"

if [[ ! -f /tmp/nomad-dtree ]]
then
  wget https://github.com/sagarrakshe/nomad-dtree/releases/download/v1.0.0/nomad-dtree -O \
    /tmp/nomad-dtree && \
    chmod +x /tmp/nomad-dtree
fi

/tmp/nomad-dtree \
  --job nginx \
  --server-addr ${NOMAD_ADDR} \
  --store consul \
  --consul-addr ${CONSUL_ADDR} \
  --consul-depfile-path nomad-dtree/dependency.json \
  --consul-jobs-path nomad-dtree/jobs
