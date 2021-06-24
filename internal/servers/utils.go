package servers

import (
	"context"
	"fmt"
	"strings"

	"github.com/kubernetes-csi/csi-lib-utils/protosanitizer"
	"gitlab.snapp.ir/snappcloud/consul-gslb-driver/pkg/gslbi"
	"google.golang.org/grpc"
	"k8s.io/klog/v2"
)

func ParseEndpoint(ep string) (string, string, error) {
	if strings.HasPrefix(strings.ToLower(ep), "unix://") || strings.HasPrefix(strings.ToLower(ep), "tcp://") {
		s := strings.SplitN(ep, "://", 2)
		if s[1] != "" {
			return s[0], s[1], nil
		}
	}
	return "", "", fmt.Errorf("invalid endpoint: %v", ep)
}

func logGRPC(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	klog.V(3).Infof("GRPC call: %s", info.FullMethod)
	klog.V(5).Infof("GRPC request: %+v", protosanitizer.StripSecrets(req))
	resp, err := handler(ctx, req)
	if err != nil {
		klog.Errorf("GRPC error: %v", err)
	} else {
		klog.V(5).Infof("GRPC response: %+v", protosanitizer.StripSecrets(resp))
	}
	return resp, err
}

func RunServers(endpoint, metricIP, metricPath string, metricPort int, ids gslbi.IdentityServer, cs gslbi.ControllerServer) {

	s := NewNonBlockingGRPCServer()
	s.Start(endpoint, metricIP, metricPath, metricPort, ids, cs)
	s.Wait()
}
