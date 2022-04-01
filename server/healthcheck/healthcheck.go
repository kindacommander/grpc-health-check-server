package healthcheck

import (
	"context"
	"sync"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/health/grpc_health_v1"
)

type HealthChecker struct {
	mu        sync.Mutex
	statusMap map[string]grpc_health_v1.HealthCheckResponse_ServingStatus
}

func (s *HealthChecker) Check(ctx context.Context, req *grpc_health_v1.HealthCheckRequest) (*grpc_health_v1.HealthCheckResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if req.Service == "" {
		return &grpc_health_v1.HealthCheckResponse{Status: grpc_health_v1.HealthCheckResponse_SERVING}, nil
	}
	if status, ok := s.statusMap[req.Service]; ok {
		return &grpc_health_v1.HealthCheckResponse{Status: status}, nil
	}

	return &grpc_health_v1.HealthCheckResponse{Status: grpc_health_v1.HealthCheckResponse_SERVICE_UNKNOWN}, nil
}

func (s *HealthChecker) Watch(req *grpc_health_v1.HealthCheckRequest, server grpc_health_v1.Health_WatchServer) error {
	logrus.Info("Serving the Watch request for health check")
	return server.Send(&grpc_health_v1.HealthCheckResponse{
		Status: grpc_health_v1.HealthCheckResponse_SERVING,
	})
}

func NewHealthChecker(services []string) *HealthChecker {
	statusMap := make(map[string]grpc_health_v1.HealthCheckResponse_ServingStatus)
	for _, s := range services {
		statusMap[s] = grpc_health_v1.HealthCheckResponse_SERVING
	}

	return &HealthChecker{
		statusMap: statusMap,
	}
}
