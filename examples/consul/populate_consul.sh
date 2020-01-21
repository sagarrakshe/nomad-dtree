#!/bin/bash

CONSUL_ADDR="http://localhost:8500"

curl \
    --request PUT \
    --data-binary @dependency.json\
    ${CONSUL_ADDR}/v1/kv/nomad-dtree/dependency.json

curl \
    --request PUT \
    --data-binary @jobs/redis.nomad \
    ${CONSUL_ADDR}/v1/kv/nomad-dtree/jobs/redis.nomad

curl \
    --request PUT \
    --data-binary @jobs/api.nomad \
    ${CONSUL_ADDR}/v1/kv/nomad-dtree/jobs/api.nomad

curl \
    --request PUT \
    --data-binary @jobs/nginx.nomad \
    ${CONSUL_ADDR}/v1/kv/nomad-dtree/jobs/nginx.nomad
