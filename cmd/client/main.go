package main

import (
	"context"
	"flag"
	"net/http"
	"os"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/snapp-cab/consul-gslb-driver/internal/connection"
	"github.com/snapp-cab/consul-gslb-driver/internal/rpc"
	"github.com/snapp-cab/consul-gslb-driver/pkg/gslbi"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"k8s.io/klog/v2"
)

const (
	// Default timeout of short GSLBI calls like GetPluginInfo
	gslbiTimeout = time.Second
)

// Command line flags
var (
	gslbiAddress = "/Users/my/gitlab/consul-gslb-driver/socket" // Address of the GSLBI driver socket.
	timeout      = flag.Duration("timeout", 15*time.Second, "Timeout for waiting for creating or deleting the gslb.")
)

func main() {
	// Connect to GSLBI.
	gslbiConn, err := connection.Connect(gslbiAddress, []grpc.DialOption{}, connection.OnConnectionLoss(connection.ExitOnConnectionLoss()))
	if err != nil {
		klog.Error(err.Error())
		os.Exit(1)
	}
	err = rpc.ProbeForever(gslbiConn, *timeout)
	if err != nil {
		klog.Error(err.Error())
		os.Exit(1)
	}

	// Find driver name.
	ctx, cancel := context.WithTimeout(context.Background(), gslbiTimeout)
	defer cancel()
	gslbiAttacher, err := rpc.GetDriverName(ctx, gslbiConn)
	if err != nil {
		klog.Error(err.Error())
		os.Exit(1)
	}
	klog.V(2).Infof("gslbi driver name: %q", gslbiAttacher)

	// Prepare http endpoint for metrics + leader election healthz
	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.Handler())
	addr := "localhost:8080"
	if addr != "" {
		go func() {
			klog.Infof("ServeMux listening at %q", addr)
			err := http.ListenAndServe(addr, mux)
			if err != nil {
				klog.Fatalf("Failed to start HTTP server at specified address (%q): %s", addr, err)
			}
		}()
	}

	h := NewAttacher(gslbiConn)
	h.Create(ctx, "hi")
	time.Sleep(100 * time.Second)
}

// Attacher implements create/delete operations against a remote gslbi driver.
type Attacher interface {
	Create(ctx context.Context, v string) (gslb string, deleted bool, err error)
	Delete(ctx context.Context) error
}

type creater struct {
	conn *grpc.ClientConn
}

// NewAttacher provides a new Attacher object.
func NewAttacher(conn *grpc.ClientConn) Attacher {
	return &creater{
		conn: conn,
	}
}

func (a *creater) Create(ctx context.Context, v string) (gslb string, deleted bool, err error) {
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

func (a *creater) Delete(ctx context.Context) error {
	client := gslbi.NewControllerClient(a.conn)

	req := gslbi.DeleteGSLBRequest{
		ServiceID: "sth",
	}

	_, err := client.DeleteGSLB(ctx, &req)
	return err
}

//// another utils

// isFinished returns true if given error represents final error of an
// operation. That means the operation has failed completely and cannot be in
// progress.  It returns false, if the error represents some transient error
// like timeout and the operation itself or previous call to the same
// operation can be actually in progress.
func isFinalError(err error) bool {
	// Source: https://github.com/grpc/grpc/blob/master/doc/statuscodes.md
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
		codes.Aborted:           // GSLBI: Operation pending for gslb
		return false
	}
	// All other errors mean that the operation (create/delete) either did not
	// even start or failed. It is for sure not in progress.
	return true
}
