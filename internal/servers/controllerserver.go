package servers

import (
	"fmt"
	"strconv"

	"gitlab.snapp.ir/snapp-cab/consul-gslb-driver/internal/consul"
	"gitlab.snapp.ir/snapp-cab/consul-gslb-driver/pkg/gslbi"
	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"k8s.io/klog/v2"
)

type controllerServer struct {
	Driver *ConsulDriver
	Consul consul.IConsul
	gslbi.UnimplementedControllerServer
}

func (cs *controllerServer) CreateGSLB(ctx context.Context, req *gslbi.CreateGSLBRequest) (*gslbi.CreateGSLBResponse, error) {
	// klog.V(4).Infof("CreateGSLB: called with args %+v", protosanitizer.StripSecrets(*req))

	// Node
	node := req.GetName()
	if len(node) == 0 {
		return nil, status.Error(codes.InvalidArgument, "[CreateGSLB] missing Gslb Name")
	}

	// serviceID
	serviceID := req.GetServiceName()
	if len(serviceID) == 0 {
		return nil, status.Error(codes.InvalidArgument, "[CreateGSLB] missing Gslb Service Name")
	}

	// Host
	address := req.GetHost()
	if len(address) == 0 {
		return nil, status.Error(codes.InvalidArgument, "[CreateGSLB] missing Gslb Host")
	}

	// Weight
	weight := req.GetWeight()
	klog.V(20).Infof("weight: %v", weight) // not implemented

	scheme := req.GetParameters()["probe_scheme"]
	probeAddress := req.GetParameters()["probe_address"]
	path := req.GetParameters()["probe_path"]
	probeFullAddress := scheme + "://" + probeAddress + path

	timeout, err := strconv.Atoi(req.GetParameters()["probe_timeout"])
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "[CreateGSLB] cannot convert probe_timeout to int")
	}
	interval, err := strconv.Atoi(req.GetParameters()["probe_interval"])
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "[CreateGSLB] cannot convert probe_interval to int")
	}

	err = cs.Consul.CreateService(node, serviceID, address, probeFullAddress, interval, timeout, make(map[string][]string))

	if err != nil {
		klog.Errorf("Failed to CreateGSLB: %v", err)
		return nil, status.Error(codes.Internal, fmt.Sprintf("CreateGSLB failed with error %v", err))
	}

	klog.V(4).Infof("CreateGSLB: Successfully created gslb %s", node)

	resp := &gslbi.CreateGSLBResponse{
		Gslb: &gslbi.Gslb{
			GslbId: "someid", //tod

		},
	}

	return resp, nil
}

func (cs *controllerServer) DeleteGSLB(ctx context.Context, req *gslbi.DeleteGSLBRequest) (*gslbi.DeleteGSLBResponse, error) {
	// klog.V(4).Infof("DeleteGSLB: called with args %+v", protosanitizer.StripSecrets(*req))

	// GLSB Delete
	gslbID := req.GetGslbId()
	if len(gslbID) == 0 {
		return nil, status.Error(codes.InvalidArgument, "DeleteGSLB GSLB ID must be provided")
	}
	err := cs.Consul.DeregService(gslbID)

	if err != nil {
		if err.Error() == "NotFound" {
			klog.V(3).Infof("Volume %s is already deleted.", gslbID)
			return &gslbi.DeleteGSLBResponse{}, nil
		}
		klog.Errorf("Failed to DeleteGSLB: %v", err)
		return nil, status.Error(codes.Internal, fmt.Sprintf("DeleteGSLB failed with error %v", err))
	}

	klog.V(4).Infof("DeleteGSLB: Successfully deleted service %s", gslbID)

	return &gslbi.DeleteGSLBResponse{}, nil
}

func (cs *controllerServer) ControllerGetGSLB(context.Context, *gslbi.ControllerGetGSLBRequest) (*gslbi.ControllerGetGSLBResponse, error) {
	return nil, status.Error(codes.Unimplemented, fmt.Sprintf("ControllerGetGSLB is not yet implemented"))
}
