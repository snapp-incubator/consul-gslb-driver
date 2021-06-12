package servers

import (
	"net"
	"os"
	"sync"

	"google.golang.org/grpc"
	"k8s.io/klog/v2"

	"github.com/snapp-cab/consul-gslb-driver/pkg/gslbi"
)

// NonBlockingGRPCServer defines Non blocking GRPC server interfaces
type NonBlockingGRPCServer interface {
	// Start services at the endpoint
	Start(endpoint string, ids gslbi.IdentityServer, cs gslbi.ControllerServer)
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

func (s *nonBlockingGRPCServer) Start(endpoint string, ids gslbi.IdentityServer, cs gslbi.ControllerServer) {

	s.wg.Add(1)

	go s.serve(endpoint, ids, cs)
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

func (s *nonBlockingGRPCServer) serve(endpoint string, ids gslbi.IdentityServer, cs gslbi.ControllerServer) {

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
		grpc.UnaryInterceptor(logGRPC),
	}
	server := grpc.NewServer(opts...)
	s.server = server

	if ids != nil {
		gslbi.RegisterIdentityServer(server, ids)
	}
	if cs != nil {
		gslbi.RegisterControllerServer(server, cs)
	}

	klog.Infof("Listening for connections on endpoint: %#v", listener.Addr())

	server.Serve(listener)

}
