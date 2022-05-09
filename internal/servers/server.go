package servers

import (
	"net"
	"net/http"
	"os"
	"strconv"
	"sync"

	"google.golang.org/grpc"
	"k8s.io/klog/v2"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/snapp-incubator/consul-gslb-driver/pkg/gslbi"
)

// NonBlockingGRPCServer defines Non blocking GRPC server interfaces
type NonBlockingGRPCServer interface {
	// Start services at the endpoint
	Start(endpoint, metricIP, metricPath string, metricPort int, ids gslbi.IdentityServer, cs gslbi.ControllerServer)
	// Waits for the service to stop
	Wait()
	// Stops the service gracefully
	Stop()
	// Stops the service forcefully
	ForceStop()
}

func NewNonBlockingGRPCServer() NonBlockingGRPCServer {
	return &nonBlockingGRPCServer{}
}

// NonBlocking server
type nonBlockingGRPCServer struct {
	wg     sync.WaitGroup
	server *grpc.Server
}

func (s *nonBlockingGRPCServer) Start(endpoint, metricIP, metricPath string, metricPort int, ids gslbi.IdentityServer, cs gslbi.ControllerServer) {

	s.wg.Add(1)

	go s.serve(endpoint, metricIP, metricPath, metricPort, ids, cs)
}

func (s *nonBlockingGRPCServer) Wait() {
	s.wg.Wait()
}

func (s *nonBlockingGRPCServer) Stop() {
	s.server.GracefulStop()
}

func (s *nonBlockingGRPCServer) ForceStop() {
	s.server.Stop()
}

func (s *nonBlockingGRPCServer) serve(endpoint, metricIP, metricPath string, metricPort int, ids gslbi.IdentityServer, cs gslbi.ControllerServer) {

	proto, addr, err := ParseEndpoint(endpoint)
	if err != nil {
		klog.Fatal(err.Error())
	}

	if proto == "unix" {
		addr = "/" + addr
		if err := os.Remove(addr); err != nil && !os.IsNotExist(err) {
			klog.Fatalf("Failed to remove %s, error: %s", addr, err.Error())
		}
	}

	listener, err := net.Listen(proto, addr)
	if err != nil {
		klog.Fatalf("Failed to listen: %v", err)
	}

	opts := []grpc.ServerOption{
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			logGRPC,
			grpc_prometheus.UnaryServerInterceptor,
		)),
	}
	server := grpc.NewServer(opts...)
	s.server = server

	if ids != nil {
		gslbi.RegisterIdentityServer(server, ids)
	}
	if cs != nil {
		gslbi.RegisterControllerServer(server, cs)
	}
	grpc_prometheus.Register(server)
	grpc_prometheus.EnableHandlingTimeHistogram()
	mux := http.NewServeMux()
	mux.Handle(metricPath, promhttp.Handler())
	metricAddr := net.JoinHostPort(metricIP, strconv.Itoa(metricPort))
	if metricAddr != "" {
		go func() {
			klog.Infof("Metrics listening at %q", metricAddr+metricPath)
			err := http.ListenAndServe(metricAddr, mux)
			if err != nil {
				klog.Fatalf("Failed to start HTTP server at specified address (%q): %s", metricAddr, err)
			}
		}()
	}

	klog.Infof("Listening for connections on endpoint: %#v", listener.Addr())

	server.Serve(listener)

}
