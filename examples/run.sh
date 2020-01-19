#!/bin/bash

make -C ../ clean
make -C ../ build

../dist/nomad-dtree \
  --job nginx \
  --server-addr http://localhost:4646/ \
  --store filesystem \
  --fs-depfile-path /home/sagar/experiments/nomad-dtree/examples/example.json \
  --fs-jobs-path /home/sagar/experiments/nomad-dtree/examples/jobs
