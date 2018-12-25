package Endpoint

import (
	"context"
	"github.com/go-kit/kit/endpoint"
	pb "google.golang.org/grpc/examples/helloworld/Proto"
)

func HelloEndpoint() endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		rep := request.(*pb.HelloRequest)
		return &pb.HelloReply{Message: "Hello " + rep.Name}, nil
	}
}
func DecodeRequest(_ context.Context, req interface{}) (interface{}, error) {
	return req, nil
}

func EncodeResponse(_ context.Context, rsp interface{}) (interface{}, error) {
	return rsp, nil
}
