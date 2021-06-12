package servers

import (
	"github.com/snapp-cab/consul-gslb-driver/internal/gslbi"
	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"k8s.io/klog/v2"
)

type identityServer struct {
	Driver *ConsulDriver
	gslbi.UnimplementedIdentityServer
}

func (ids *identityServer) GetPluginInfo(ctx context.Context, req *gslbi.GetPluginInfoRequest) (*gslbi.GetPluginInfoResponse, error) {
	klog.V(5).Infof("Using default GetPluginInfo")

	if ids.Driver.name == "" {
		return nil, status.Error(codes.Unavailable, "Driver name not configured")
	}

	if ids.Driver.fqVersion == "" {
		return nil, status.Error(codes.Unavailable, "Driver is missing version")
	}

	return &gslbi.GetPluginInfoResponse{
		Name:          ids.Driver.name,
		VendorVersion: ids.Driver.fqVersion,
	}, nil
}

func (ids *identityServer) Probe(ctx context.Context, req *gslbi.ProbeRequest) (*gslbi.ProbeResponse, error) {
	// oProvider, err := openstack.GetOpenStackProvider()
	// if err != nil {
	// 	klog.Errorf("Failed to GetOpenStackProvider: %v", err)
	// 	return nil, status.Error(codes.FailedPrecondition, "Failed to retrieve openstack provider")
	// }
	// if err := oProvider.CheckBlockStorageAPI(); err != nil {
	// 	klog.Errorf("Failed to query blockstorage API: %v", err)
	// 	return nil, status.Error(codes.FailedPrecondition, "Failed to communicate with OpenStack BlockStorage API")
	// }
	return &gslbi.ProbeResponse{}, nil
}

func (ids *identityServer) GetPluginCapabilities(ctx context.Context, req *gslbi.GetPluginCapabilitiesRequest) (*gslbi.GetPluginCapabilitiesResponse, error) {
	klog.V(5).Infof("GetPluginCapabilities called with req %+v", req)
	return &gslbi.GetPluginCapabilitiesResponse{
		Capabilities: []*gslbi.PluginCapability{
			{
				Type: &gslbi.PluginCapability_Service_{
					Service: &gslbi.PluginCapability_Service{
						Type: gslbi.PluginCapability_Service_CONTROLLER_SERVICE,
					},
				},
			},
			{
				Type: &gslbi.PluginCapability_HealthCheck_{
					HealthCheck: &gslbi.PluginCapability_HealthCheck{
						Type: gslbi.PluginCapability_HealthCheck_HTTP,
					},
				},
			},
			{
				Type: &gslbi.PluginCapability_HealthCheck_{
					HealthCheck: &gslbi.PluginCapability_HealthCheck{
						Type: gslbi.PluginCapability_HealthCheck_TCP,
					},
				},
			},
		},
	}, nil
}
