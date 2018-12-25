

//protoc -I ../Proto --go_out=plugins=grpc:../Proto ../Proto/Proto.proto

package main

import (
	"fmt"
	grpc_transport "github.com/go-kit/kit/transport/grpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/examples/helloworld/Endpoint"
	"google.golang.org/grpc/examples/helloworld/Service"
	pb "google.golang.org/grpc/examples/helloworld/Proto"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
)

const (
	port = ":50051"
)


func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	helloServer := new(Service.Server)
	SayHelloHandler := grpc_transport.NewServer(
		Endpoint.HelloEndpoint(),
		Endpoint.DecodeRequest,
		Endpoint.EncodeResponse,
	)
	helloServer.HelloHandler=SayHelloHandler

	s := grpc.NewServer()
	pb.RegisterGreeterServer(s, helloServer)
	reflection.Register(s)
	fmt.Print("server is runing")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

}

