rm -f ./phpfestproto/*.go
protoc --go_out=:. --go-grpc_out=:. phpfestproto/*.proto