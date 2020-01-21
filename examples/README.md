Prerequisite
============

To run the following examples you need to have nomad and consul running.
Assuming the nomad and consul is running on local, if their address is different
just the address in the `run.sh` and run the script.

For consul example, run `populate.sh` script before running `run.sh`. The former
script pushes the nomad job files and dependency json on consul.
