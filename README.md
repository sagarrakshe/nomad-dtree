
nomad-dtree
===========

Nomad doesn't have its own dependency management plugin between the jobs. 
Especially with rise in the microservices pattern, it becomes difficult to maintain order of deployment using Nomad.
The order of deployment of the jobs has to be handled separately, using some external tool like Rundeck.
nomad-dtree is an attempt to address this problem.

Getting Started
===============

nomad-dtree requires a dependency and all nomad job files. These dependency file and nomad jobs can be stored on filesystem or consul.
The dependencies are defined in JSON format. Check [Sample dependecy json ](https://github.com/sagarrakshe/nomad-dtree/blob/master/examples/example.json)

```
{
  "dependencies": {
    "nginx": {
      "wait_cond": 5,
        "pre": {
	  "job": "api"
	  }
	},
    "api": {
      "wait_cond": 7,
      "pre": {
        "job": "redis"
      }
    },
    "redis": {
      "wait_cond": 5
    }
  }
}
```
The `wait_cond` denotes the wait time before running the next job.
```
nomad-dtree \
  --job <job_name> \
  --server-addr <nomad-address> \
  --store <filesystem or consul> \
  --fs-depfile-path <path to dependency json file> \
  --fs-jobs-path <path to directory containing jobs>
```

Example
=======

Run nomad on your local machine. If the nomad address is different from http://localhost:4646 change the address in `examples/run.sh` accordingly.
Download the nomad-dtree binary from [here](https://github.com/sagarrakshe/nomad-dtree/releases/download/v1.0.0/nomad-dtree).

```
$ wget https://github.com/sagarrakshe/nomad-dtree/releases/download/v1.0.0/nomad-dtree -O /tmp/nomad-dtree && chmod +x /tmp/nomad-dtree
$ cd examples
$ bash run.sh
```
