package rpc

import (
	"context"
	"fmt"
	"time"

	"github.com/snapp-cab/consul-gslb-driver/internal/gslbi"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"k8s.io/klog/v2"
)

const (
	// Interval of trying to call Probe() until it succeeds
	probeInterval = 1 * time.Second
)

// ProbeForever calls Probe() of a CSI driver and waits until the driver becomes ready.
// Any error other than timeout is returned.
func ProbeForever(conn *grpc.ClientConn, singleProbeTimeout time.Duration) error {
	for {
		klog.Info("Probing CSI driver for readiness")
		ready, err := probeOnce(conn, singleProbeTimeout)
		if err != nil {
			st, ok := status.FromError(err)
			if !ok {
				// This is not gRPC error. The probe must have failed before gRPC
				// method was called, otherwise we would get gRPC error.
				return fmt.Errorf("CSI driver probe failed: %s", err)
			}
			if st.Code() != codes.DeadlineExceeded {
				return fmt.Errorf("CSI driver probe failed: %s", err)
			}
			// Timeout -> driver is not ready. Fall through to sleep() below.
			klog.Warning("CSI driver probe timed out")
		} else {
			if ready {
				return nil
			}
			klog.Warning("CSI driver is not ready")
		}
		// Timeout was returned or driver is not ready.
		time.Sleep(probeInterval)
	}
}

// probeOnce is a helper to simplify defer cancel()
func probeOnce(conn *grpc.ClientConn, timeout time.Duration) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	return Probe(ctx, conn)
}

// Probe calls driver Probe() just once and returns its result without any processing.
func Probe(ctx context.Context, conn *grpc.ClientConn) (ready bool, err error) {
	client := gslbi.NewIdentityClient(conn)

	req := gslbi.ProbeRequest{}
	rsp, err := client.Probe(ctx, &req)

	if err != nil {
		return false, err
	}

	r := rsp.GetReady()
	if r == nil {
		// "If not present, the caller SHALL assume that the plugin is in a ready state"
		return true, nil
	}
	return r.GetValue(), nil
}

// GetDriverName returns name of CSI driver.
func GetDriverName(ctx context.Context, conn *grpc.ClientConn) (string, error) {
	client := gslbi.NewIdentityClient(conn)

	req := gslbi.GetPluginInfoRequest{}
	rsp, err := client.GetPluginInfo(ctx, &req)
	if err != nil {
		return "", err
	}
	name := rsp.GetName()
	if name == "" {
		return "", fmt.Errorf("driver name is empty")
	}
	return name, nil
}
