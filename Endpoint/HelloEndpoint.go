package Endpoint

import (
	"context"
	"github.com/go-kit/kit/circuitbreaker"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/ratelimit"
	"github.com/go-kit/kit/tracing/zipkin"
	"github.com/sony/gobreaker"
	"golang.org/x/time/rate"
	gozipkin "github.com/openzipkin/zipkin-go"
	pb "google.golang.org/grpc/examples/helloworld/Proto"
)

func Hellopoint(limit *rate.Limiter,zipkinTracer *gozipkin.Tracer) endpoint.Endpoint{
	var HelloPoint endpoint.Endpoint
	{
		HelloPoint = HelloEndpoint()
		//rate.NewLimiter(rate.Every(time.Second), 1)
		HelloPoint = ratelimit.NewErroringLimiter(limit)(HelloPoint)
		HelloPoint = circuitbreaker.Gobreaker(gobreaker.NewCircuitBreaker(gobreaker.Settings{}))(HelloPoint)
		//HelloPoint = opentracing.TraceServer()(HelloPoint)
		HelloPoint = zipkin.TraceEndpoint(zipkinTracer,"hello")(HelloPoint)
	}

	return HelloPoint
}

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
