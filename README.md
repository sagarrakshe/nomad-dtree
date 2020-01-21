
nomad-dtree
===========
Nomad doesn't have its own dependency management plugin between the jobs. 
Especially with rise in the microservices pattern, it becomes difficult to maintain order of deployment using Nomad.
For eg. Database service must be up before the api service is started.
The order of deployment of the jobs has to be handled separately, using some external tool like Rundeck.

Nomad also lacks feature or post/pre hook. (Check issue: [https://github.com/hashicorp/nomad/issues/419](https://github.com/hashicorp/nomad/issues/419) and [https://github.com/hashicorp/nomad/issues/2851](https://github.com/hashicorp/nomad/issues/2851)).
Many of services a pre-setup step or a post setup.
For example, when we run a postgres or any db container, we need to create a user, db and grant permissions.
Another example, on bringing up a rabbitmq cluster we might need to create vhost, user etc.

**nomad-dtree** is an attempt to address these problems.

nomad-dtree expects everything to be a nomad job. Even a post/pre hook should be a nomad jobs.
It also expects a JSON where the dependency is mentioned.


How it works
=============

Consider a small service which has three components:

- nginx
- api
- postgres

The order of dependency is same as above. The postgres service has a post hook, a
small step of creating a user and db which is used by the api service. So we have
nomad jobs:

- nginx.nomad
- api.nomad
- postgres.nomad
- postgres_setup.nomad

The dependency file for the above servcies is as follows:

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
        "job": "postgres"
      }
    },
    "postgres": {
      "post": {
        "job": "postgres_setup",
        "wait_cond": 3
      }
      "wait_cond": 5
    }
  }
}
```

Now, each block in the `dependencies` json denotes a nomad job. It has three
parts, consider `api`:

```
    "api": {
      "wait_cond": 7,
      "pre": {
        "job": "postgres"
      }
    },
```

- **wait_cond** (required) - It the wait_time, time taken for the api service to be healthy. 
- **pre** (optional) - Pre-requisite job that has to be done before running the api, in this case postgres should be running before running postgres.
- **post**(optional) - Post hook, job that needs to be run after api is running (and after waiting for `wait_cond` seconds).

Now, when we run the following command:

```
nomad-dtree \
  --job nginx \  (note we are submitting nginx job here)
  --server-addr <nomad-address> \
  --store filesystem \
  --fs-depfile-path <path to dependency json file> \
  --fs-jobs-path <path to directory containing all nomad jobs>
```

The order of execution will be as follows:

```
postgres.nomad -> postgres_setup.nomad -> api.nomad -> nginx.nomad
```
For more examples check  [Examples](!https://github.com/sagarrakshe/nomad-dtree/tree/master/examples)
