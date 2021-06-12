package consul

import (
	"fmt"
	"time"

	consulapi "github.com/hashicorp/consul/api"
)

func (c *Consul) DeregService(serviceID string) error {
	dereg := &consulapi.CatalogDeregistration{
		Node: serviceID,
	}
	_, err := c.catalog.Deregister(dereg, &consulapi.WriteOptions{})
	if err != nil {
		return fmt.Errorf("failed to deregister consul service: %w", err)
	}
	return nil
}

func (c *Consul) CreateService(node, serviceID, address, probeAddress string, intervalDuration, timeoutDuration int, headers map[string][]string) error {

	reg := &consulapi.CatalogRegistration{
		Node:    node,
		Address: address,
		NodeMeta: map[string]string{
			"external-node":  "true",
			"external-probe": "false",
		},
		Service: &consulapi.AgentService{
			ID:      serviceID,
			Service: serviceID,
		},
		Checks: consulapi.HealthChecks{
			&consulapi.HealthCheck{
				Name:   "http-check",
				Status: "passing",
				Definition: consulapi.HealthCheckDefinition{
					HTTP:             probeAddress,
					IntervalDuration: time.Duration(intervalDuration) * time.Second,
					TimeoutDuration:  time.Duration(timeoutDuration) * time.Second,
					Header:           headers,
				},
			},
		},
	}
	_, err := c.catalog.Register(reg, &consulapi.WriteOptions{})
	if err != nil {
		return fmt.Errorf("failed to register consul service: %w", err)
	}
	return nil
}
