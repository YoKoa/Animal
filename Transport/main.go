

//protoc -I ../Proto --go_out=plugins=grpc:../Proto ../Proto/Proto.proto

package main

import (
	"Animal/Endpoint"
	"fmt"
	grpc_transport "github.com/go-kit/kit/transport/grpc"
	"github.com/openzipkin/zipkin-go"
	zipkinhttp "github.com/openzipkin/zipkin-go/reporter/http"
	"golang.org/x/time/rate"
	"google.golang.org/grpc"
	pb "google.golang.org/grpc/examples/helloworld/Proto"
	"google.golang.org/grpc/examples/helloworld/Service"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
	"os"
	"time"
)

const (
	port = ":50051"
	endpointURL = "http://localhost:9411/api/v2/spans"
)
//var fs = flag.NewFlagSet("hello", flag.ExitOnError)
//var (
//	//debugAddr      = fs.String("debug.addr", ":8080", "Debug and metrics listen address")
//	//httpAddr       = fs.String("http-addr", ":8081", "HTTP listen address")
//	//grpcAddr       = fs.String("grpc-addr", ":8082", "gRPC listen address")
//	//thriftAddr     = fs.String("thrift-addr", ":8083", "Thrift listen address")
//	//jsonRPCAddr    = fs.String("jsonrpc-addr", ":8084", "JSON RPC listen address")
//	//thriftProtocol = fs.String("thrift-protocol", "binary", "binary, compact, json, simplejson")
//	//thriftBuffer   = fs.Int("thrift-buffer", 0, "0 for unbuffered")
//	//thriftFramed   = fs.Bool("thrift-framed", false, "true to enable framing")
//	zipkinV2URL    = fs.String("http://localhost:9411/api/v2/spans", "", "Enable Zipkin v2 tracing (zipkin-go) using a Reporter URL e.g. http://localhost:9411/api/v2/spans")
//	//zipkinV1URL    = fs.String("http://localhost:9411/api/v1/spans", "", "Enable Zipkin v1 tracing (zipkin-go-opentracing) using a collector URL e.g. http://localhost:9411/api/v1/spans")
//	//lightstepToken = fs.String("lightstep-token", "", "Enable LightStep tracing via a LightStep access token")
//	//appdashAddr    = fs.String("appdash-addr", "", "Enable Appdash tracing via an Appdash server host:port")
//)
var limit = rate.NewLimiter(rate.Every(time.Second), 1)



func main() {

	var zipkinTracer *zipkin.Tracer
	{
		var(
			err  error
			hostPort= "localhost:0"
			serviceName= "hello"
			reporter= zipkinhttp.NewReporter(endpointURL)
		)
		defer reporter.Close()
		sampler, err := zipkin.NewCountingSampler(1)
		zEP, _ := zipkin.NewEndpoint(serviceName, hostPort)
		zipkinTracer, err = zipkin.NewTracer(
			reporter,
			zipkin.WithSampler(sampler),
			zipkin.WithLocalEndpoint(zEP),
		)
		if err != nil {
			log.Fatal("err", err)
			os.Exit(1)
		}
	}

	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	helloServer := new(Service.Server)
	SayHelloHandler := grpc_transport.NewServer(
		Endpoint.Hellopoint(limit,zipkinTracer),
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

