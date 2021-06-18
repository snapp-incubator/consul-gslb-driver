package connection

import (
	"context"
	"errors"
	"io/ioutil"
	"net"
	"strings"
	"time"

	"github.com/kubernetes-csi/csi-lib-utils/protosanitizer"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"k8s.io/klog/v2"
)

const (
	// Interval of logging connection errors
	connectionLoggingInterval = 10 * time.Second
)
const terminationLogPath = "/dev/termination-log"

// connect is the internal implementation of Connect. It has more options to enable testing.
func Connect(
	address string,
	dialOptions []grpc.DialOption, connectOptions ...Option) (*grpc.ClientConn, error) {
	var o options
	for _, option := range connectOptions {
		option(&o)
	}

	dialOptions = append(dialOptions,
		grpc.WithInsecure(),                   // Don't use TLS, it's usually local Unix domain socket in a container.
		grpc.WithBackoffMaxDelay(time.Second), // Retry every second after failure.
		grpc.WithBlock(),                      // Block until connection succeeds.
		grpc.WithChainUnaryInterceptor(
			LogGRPC, // Log all messages.
		),
	)
	unixPrefix := "unix://"
	if strings.HasPrefix(address, "/") {
		// It looks like filesystem path.
		address = unixPrefix + address
	}

	if strings.HasPrefix(address, unixPrefix) {
		// state variables for the custom dialer
		haveConnected := false
		lostConnection := false
		reconnect := true

		dialOptions = append(dialOptions, grpc.WithDialer(func(addr string, timeout time.Duration) (net.Conn, error) {
			if haveConnected && !lostConnection {
				// We have detected a loss of connection for the first time. Decide what to do...
				// Record this once. TODO (?): log at regular time intervals.
				klog.Errorf("Lost connection to %s.", address)
				// Inform caller and let it decide? Default is to reconnect.
				if o.reconnect != nil {
					reconnect = o.reconnect()
				}
				lostConnection = true
			}
			if !reconnect {
				return nil, errors.New("connection lost, reconnecting disabled")
			}
			conn, err := net.DialTimeout("unix", address[len(unixPrefix):], timeout)
			if err == nil {
				// Connection reestablished.
				haveConnected = true
				lostConnection = false
			}
			return conn, err
		}))
	} else if o.reconnect != nil {
		return nil, errors.New("OnConnectionLoss callback only supported for unix:// addresses")
	}

	klog.V(5).Infof("Connecting to %s", address)

	// Connect in background.
	var conn *grpc.ClientConn
	var err error
	ready := make(chan bool)
	go func() {
		conn, err = grpc.Dial(address, dialOptions...)
		close(ready)
	}()

	// Log error every connectionLoggingInterval
	ticker := time.NewTicker(connectionLoggingInterval)
	defer ticker.Stop()

	// Wait until Dial() succeeds.
	for {
		select {
		case <-ticker.C:
			klog.Warningf("Still connecting to %s", address)

		case <-ready:
			return conn, err
		}
	}
}

// Option is the type of all optional parameters for Connect.
type Option func(o *options)
type options struct {
	reconnect func() bool
}

// LogGRPC is gPRC unary interceptor for logging of GSLB messages at level 5. It removes any secrets from the message.
func LogGRPC(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	klog.V(5).Infof("GRPC call: %s", method)
	klog.V(5).Infof("GRPC request: %s", protosanitizer.StripSecrets(req))
	err := invoker(ctx, method, req, reply, cc, opts...)
	klog.V(5).Infof("GRPC response: %s", protosanitizer.StripSecrets(reply))
	klog.V(5).Infof("GRPC error: %v", err)
	return err
}

// OnConnectionLoss registers a callback that will be invoked when the
// connection got lost. If that callback returns true, the connection
// is reestablished. Otherwise the connection is left as it is and
// all future gRPC calls using it will fail with status.Unavailable.
func OnConnectionLoss(reconnect func() bool) Option {
	return func(o *options) {
		o.reconnect = reconnect
	}
}

// ExitOnConnectionLoss returns callback for OnConnectionLoss() that writes
// an error to /dev/termination-log and exits.
func ExitOnConnectionLoss() func() bool {
	return func() bool {
		terminationMsg := "Lost connection to GSLB driver, exiting"
		if err := ioutil.WriteFile(terminationLogPath, []byte(terminationMsg), 0644); err != nil {
			klog.Errorf("%s: %s", terminationLogPath, err)
		}
		klog.Fatalf(terminationMsg)
		return false
	}
}

// IsFinished returns true if given error represents final error of an
// operation. That means the operation has failed completely and cannot be in
// progress.  It returns false, if the error represents some transient error
// like timeout and the operation itself or previous call to the same
// operation can be actually in progress.
func IsFinalError(err error) bool {
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
