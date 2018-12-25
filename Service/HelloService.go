package Service
import (
	"context"
	grpc_transport "github.com/go-kit/kit/transport/grpc"
	pb "google.golang.org/grpc/examples/helloworld/Proto"
)
type Server struct{
	HelloHandler  grpc_transport.Handler
}

func (s *Server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	_,response,err :=s.HelloHandler.ServeGRPC(ctx,in)
	if err != nil{
		return nil,err
	}
	return response.(*pb.HelloReply),nil
}
