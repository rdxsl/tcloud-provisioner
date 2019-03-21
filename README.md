# tcloud-provisioner

## Build Instructions

Run `make build` to compile the binary tool.

Run `make` to run unit tests and compile the binary tool.

## Install Instructions

Run `go install github.com/rdxsl/tcloud-provisioner` and make sure `$GOPATH/bin` is on your `$PATH`

## MySQL Idempotency

The  `instancename` in `mysql.json` is Idempotency, this means if a MySQL instance in tcloud has the same `instancename`, this program will not create a duplicate instance with the same name.

When you delete the mysql instance in tclound, please make sure you release the deleted instance in the `recyle-bin`. 
