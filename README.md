Demo project with gRPC and protobuf especially for PHPFest conference

## Environment installations

### 1. protoc

see http://google.github.io/proto-lens/installing-protoc.html

### 2. protoc-gen-go

see https://grpc.io/docs/languages/go/quickstart/

### 3. protoc-gen-go-grpc

see https://godoc.org/google.golang.org/grpc/cmd/protoc-gen-go-grpc

### 4. Binaries

move binary files (protoc-gen-go, protoc-gen-go-grpc) to the /usr/bin or to the /usr/local/bin library 

## Project installation

```
cd $GOPATH/src
git clone https://github.com/antondzhukov/phpfest-go.git
cd phpfest-go
make protofiles
make install
```

## Use

Via tcp
```
./phpfestgo --listen-proto=tcp --listen-address=127.0.0.1:8282
```

Via unix socket
```
./phpfestgo --listen-proto=unix --listen-address=phpfestgo.sock
```