package consul

import (
	consulapi "github.com/hashicorp/consul/api"
)

var ConsulInstance IConsul = nil

// var configFiles = []string{"/etc/cloud.conf"}
// var cfg Config

type IConsul interface {
	CreateService(node, serviceID, address, probeAddress string, intervalDuration, timeoutDuration int, headers map[string][]string) error
	DeregService(serviceID string) error
}

type Consul struct {
	client  *consulapi.Client
	catalog *consulapi.Catalog
}

// CreateConsulProvider creates Consul Instance
func CreateConsulProvider() (IConsul, error) {
	// Get config from file
	// cfg, err := GetConfigFromFiles(configFiles)
	// if err != nil {
	// 	klog.Errorf("GetConfigFromFiles %s failed with error: %v", configFiles, err)
	// 	return nil, err
	// }
	// logcfg(cfg)

	config := &consulapi.Config{
		Address:    "consul.apps.private.okd4.teh-1.snappcloud.io",
		Scheme:     "http",
		Datacenter: "teh1",
	}
	client, err := consulapi.NewClient(config)
	if err != nil {
		return nil, err
	}

	// Init Consul
	ConsulInstance = &Consul{
		client:  client,
		catalog: client.Catalog(),
	}

	return ConsulInstance, nil
}

// GetConsul returns Consul Instance
func GetConsul() (IConsul, error) {
	if ConsulInstance != nil {
		return ConsulInstance, nil
	}
	var err error
	ConsulInstance, err = CreateConsulProvider()
	if err != nil {
		return nil, err
	}

	return ConsulInstance, nil
}
