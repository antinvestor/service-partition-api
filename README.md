# service-partition-api

A repository for the partition service api that ants use to group into logical and physical boundaries while investing.

### How do I update the definitions? ###

* The api definition is defined in the proto file partition.proto
* To update the proto service you need to run the commands :


    `protoc --proto_path=../apis --proto_path=./v1 --go_out=./ --validate_out=lang=go:. partition.proto`

    `protoc --proto_path=../apis --proto_path=./v1  partition.proto --go-grpc_out=./ `
    
    `mockgen -source=partition_grpc.pb.go -self_package=github.com/antinvestor/service-partition-api -package=partition_v1 -destination=partition_grpc_mock.go`


  with that in place update the implementation appropriately
