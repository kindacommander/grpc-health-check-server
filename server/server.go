package main

import (
	"flag"
	"fmt"
	"grpc-health-check/proto"
	"grpc-health-check/server/healthcheck"
	"net"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
)

type server struct{}

func init() {
	logrus.SetFormatter(&logrus.JSONFormatter{
		PrettyPrint: true,
	})
}

func (s *server) Hello(helloReq *proto.HelloRequest, srv proto.GreetService_HelloServer) error {
	logrus.Infof("Server received an rpc request with the following parameter %v", helloReq.Hello)

	for i := 0; i <= 10; i++ {
		resp := &proto.HelloResponse{
			Greet: fmt.Sprintf("Hello %s for %d time", helloReq.Hello, i),
		}
		srv.SendMsg(resp)
	}
	return nil
}

func main() {
	port := flag.String("p", "50051", "grpc server port")
	flag.Parse()

	serverAddress := (":" + *port)

	listenAddr, err := net.Listen("tcp", serverAddress)
	if err != nil {
		logrus.Fatalf("Error while starting the listening service %v", err.Error())
	}

	grpcServer := grpc.NewServer()
	proto.RegisterGreetServiceServer(grpcServer, &server{})

	reflection.Register(grpcServer)

	serviceInfo := grpcServer.GetServiceInfo()
	services := make([]string, 0, len(serviceInfo))
	for k := range serviceInfo {
		services = append(services, k)
	}

	healthService := healthcheck.NewHealthChecker(services)
	grpc_health_v1.RegisterHealthServer(grpcServer, healthService)

	logrus.Infof("Server started on %s", serverAddress)
	if err = grpcServer.Serve(listenAddr); err != nil {
		logrus.Fatalf("Error while starting the gRPC server on the %s listen address %v", listenAddr, err.Error())
	}
}
