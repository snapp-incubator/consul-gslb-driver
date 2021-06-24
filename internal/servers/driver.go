package servers

import (
	"gitlab.snapp.ir/snappcloud/consul-gslb-driver/internal/consul"
	"k8s.io/klog/v2"
)

const (
	driverName = "consul.gslbi.snappcloud.io"
)

var (
	specVersion = "1.0.0"
	Version     = "1.0.1"
)

type ConsulDriver struct {
	name       string
	fqVersion  string
	endpoint   string
	metricIP   string
	metricPort int
	metricPath string
	ids        *identityServer
	cs         *controllerServer
}

func NewDriver(endpoint, metricIP, metricPath string, metricPort int) *ConsulDriver {

	d := &ConsulDriver{}
	d.name = driverName
	d.fqVersion = Version
	d.endpoint = endpoint
	d.metricIP = metricIP
	d.metricPort = metricPort
	d.metricPath = metricPath

	klog.Info("Driver: ", d.name)
	klog.Info("Driver version: ", d.fqVersion)
	klog.Info("GSLBI Spec version: ", specVersion)

	return d
}

func (d *ConsulDriver) SetupDriver(consul consul.IConsul) {
	d.ids = NewIdentityServer(d)
	d.cs = NewControllerServer(d, consul)
}

func (d *ConsulDriver) Run() {
	RunServers(d.endpoint, d.metricIP, d.metricPath, d.metricPort, d.ids, d.cs)
}

func NewIdentityServer(d *ConsulDriver) *identityServer {
	return &identityServer{
		Driver: d,
	}
}

func NewControllerServer(d *ConsulDriver, consul consul.IConsul) *controllerServer {
	return &controllerServer{
		Driver: d,
		Consul: consul,
	}
}
