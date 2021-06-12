package servers

import (
	"fmt"

	"github.com/snapp-cab/consul-gslb-driver/internal/consul"
	"github.com/snapp-cab/consul-gslb-driver/pkg/gslbi"
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

	// Volume Name
	gslbName := req.GetName()

	if len(gslbName) == 0 {
		return nil, status.Error(codes.InvalidArgument, "[CreateGSLB] missing Gslb Name")
	}

	// vol, err := consul.CreateGSLB(volName, volSizeGB, volType, volAvailability, snapshotID, sourcevolID, &properties)

	// if err != nil {
	// 	klog.Errorf("Failed to CreateGSLB: %v", err)
	// 	return nil, status.Error(codes.Internal, fmt.Sprintf("CreateGSLB failed with error %v", err))

	// }

	klog.V(4).Infof("CreateGSLB: Successfully created volume %s", gslbName)

	return getCreateGSLBResponse(gslbName), nil
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

func getCreateGSLBResponse(vol string) *gslbi.CreateGSLBResponse {

	resp := &gslbi.CreateGSLBResponse{
		GSLB: vol,
	}

	return resp

}
