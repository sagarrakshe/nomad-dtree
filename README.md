nomad-dtree
===========

Tool to handle dependencies of nomad jobs.

To run the tool run following commands:

```
$ make build

$ ./dist/nomad-dtree --help

usage: nomad-dtree [<flags>]

Flags:
      --help                 Show context-sensitive help (also try --help-long and --help-man).
  -j, --job=JOB              Job
      --root-ca-file=ROOT-CA-FILE
                             RootCA File
      --cert-file=CERT-FILE  Cert File
      --key-file=KEY-FILE    Key File
      --dependency-file=DEPENDENCY-FILE
                             Dependency File
      --jobs-path=JOBS-PATH  Path to jobs location
  -s, --server-addr="http://127.0.0.1:4646"
                             Server Addr


$ ./dist/nomad-dtree -j <job_name> \
		--root-ca-file ./dist/cert-chain.pem \
		--cert-file ./dist/client.pem \
		--key-file ./dist/client-key.pem \
		--dependency-file ./dist/dependency.json \
		--server-addr https://nomaddashboard:4646/ \
		--jobs-path ./dist/nomad_jobs
```
