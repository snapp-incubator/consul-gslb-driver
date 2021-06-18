package servers

import (
	"gitlab.com/snapp-cab/consul-gslb-driver/internal/consul"
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
	name         string
	fqVersion    string
	endpoint     string
	datacenter   string
	consulConfig string
	ids          *identityServer
	cs           *controllerServer
}

func NewDriver(endpoint, datacenter string) *ConsulDriver {

	d := &ConsulDriver{}
	d.name = driverName
	d.fqVersion = Version
	d.endpoint = endpoint
	d.datacenter = datacenter

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
	RunServers(d.endpoint, d.ids, d.cs)
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
