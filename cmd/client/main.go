package main

import (
	"context"
	"flag"
	"os"
	"time"

	"k8s.io/klog/v2"

	"github.com/kubernetes-csi/csi-lib-utils/connection"
	"github.com/kubernetes-csi/csi-lib-utils/metrics"
	"github.com/kubernetes-csi/csi-lib-utils/rpc"

	gslbi "github.com/snapp-cab/consul-gslb-driver/internal/gslbi"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (

	// Default timeout of short CSI calls like GetPluginInfo
	csiTimeout = time.Second
)

// Command line flags
var (
	csiAddress = flag.String("csi-address", "unix://Users/my/gitlab/consul-gslb-driver/socket", "Address of the CSI driver socket.")
	timeout    = flag.Duration("timeout", 15*time.Second, "Timeout for waiting for attaching or detaching the volume.")
)

var (
	version = "unknown"
)

type leaderElection interface {
	Run() error
	WithNamespace(namespace string)
}

func main() {

	metricsManager := metrics.NewCSIMetricsManager("" /* driverName */)
	// Connect to CSI.
	klog.Warning("Here")
	csiConn, err := connection.Connect(*csiAddress, metricsManager, connection.OnConnectionLoss(connection.ExitOnConnectionLoss()))
	if err != nil {
		klog.Error(err.Error())
		os.Exit(1)
	}
	err = rpc.ProbeForever(csiConn, *timeout)
	if err != nil {
		klog.Error(err.Error())
		os.Exit(1)
	}

	ctx, cancel := context.WithTimeout(context.Background(), csiTimeout)
	defer cancel()
	h := NewAttacher(csiConn)
	h.Create(ctx, "hi")
}

// Attacher implements attach/detach operations against a remote CSI driver.
type Attacher interface {
	Create(ctx context.Context, v string) (gslb string, detached bool, err error)
	Delete(ctx context.Context) error
}

type attacher struct {
	conn *grpc.ClientConn
}

// NewAttacher provides a new Attacher object.
func NewAttacher(conn *grpc.ClientConn) Attacher {
	return &attacher{
		conn: conn,
	}
}

func (a *attacher) Create(ctx context.Context, v string) (gslb string, detached bool, err error) {
	client := gslbi.NewControllerClient(a.conn)

	req := gslbi.CreateGSLBRequest{
		Name: "sth",
	}

	rsp, err := client.CreateGSLB(ctx, &req)
	if err != nil {
		return "", isFinalError(err), err
	}
	return rsp.GSLB, false, nil
}

func (a *attacher) Delete(ctx context.Context) error {
	client := gslbi.NewControllerClient(a.conn)

	req := gslbi.DeleteGSLBRequest{
		ServiceID: "sth",
	}

	_, err := client.DeleteGSLB(ctx, &req)
	return err
}

func logGRPC(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	klog.V(5).Infof("GRPC call: %s", method)
	// klog.V(5).Infof("GRPC request: %s", protosanitizer.StripSecrets(req))
	err := invoker(ctx, method, req, reply, cc, opts...)
	// klog.V(5).Infof("GRPC response: %s", protosanitizer.StripSecrets(reply))
	klog.V(5).Infof("GRPC error: %v", err)
	return err
}

// isFinished returns true if given error represents final error of an
// operation. That means the operation has failed completely and cannot be in
// progress.  It returns false, if the error represents some transient error
// like timeout and the operation itself or previous call to the same
// operation can be actually in progress.
func isFinalError(err error) bool {
	// Sources:
	// https://github.com/grpc/grpc/blob/master/doc/statuscodes.md
	// https://github.com/container-storage-interface/spec/blob/master/spec.md
	st, ok := status.FromError(err)
	if !ok {
		// This is not gRPC error. The operation must have failed before gRPC
		// method was called, otherwise we would get gRPC error.
		return false
	}
	switch st.Code() {
	case codes.Canceled, // gRPC: Client Application cancelled the request
		codes.DeadlineExceeded,  // gRPC: Timeout
		codes.Unavailable,       // gRPC: Server shutting down, TCP connection broken - previous Attach() or Detach() may be still in progress.
		codes.ResourceExhausted, // gRPC: Server temporarily out of resources - previous Attach() or Detach() may be still in progress.
		codes.Aborted:           // CSI: Operation pending for volume
		return false
	}
	// All other errors mean that the operation (attach/detach) either did not
	// even start or failed. It is for sure not in progress.
	return true
}
