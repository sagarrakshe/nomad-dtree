#!/bin/bash

/tmp/nomad-dtree \
  --job nginx \
  --server-addr http://localhost:4646/ \
  --store filesystem \
  --fs-depfile-path example.json \
  --fs-jobs-path jobs
