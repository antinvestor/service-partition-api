# service-partition-api

A repository for the partition service api that ants use to group into logical and physical boundaries while investing.

### How do I update the definitions? ###

* The api definition is defined in the proto file partition.proto
* To update the proto service you need to run the commands :

  `protoc -I ../common/service common/common.proto partition/v1/partition.proto --go_out=./`
  `protoc -I ../common/service common/common.proto partition/v1/partition.proto --go-grpc_out=./`

  with that in place update the implementation appropriately
