package env

import (
	"errors"
	"net"
	"os"

	"github.com/kms-qwe/chat-server/internal/config"
)

const (
	grpcHostEnvName = "GRPC_HOST"
	gprcPortEnvName = "GRPC_PORT"
)

type grpcConfig struct {
	host string
	port string
}

// Address provides grpc serv address
func (g *grpcConfig) Address() string {
	return net.JoinHostPort(g.host, g.port)
}

// NewGRPCConfig creates a new gRPC configuration based on environment variables.
func NewGRPCConfig() (config.GRPCConfig, error) {
	host := os.Getenv(grpcHostEnvName)
	if len(host) == 0 {
		return nil, errors.New("grpc host not found")
	}

	port := os.Getenv(gprcPortEnvName)
	if len(port) == 0 {
		return nil, errors.New("grpc port not found")
	}

	return &grpcConfig{
		host: host,
		port: port,
	}, nil
}
