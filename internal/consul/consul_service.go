package consul

import (
	"fmt"

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

func (c *Consul) CreateService(serviceID string) error {
	dereg := &consulapi.CatalogDeregistration{
		Node: serviceID,
	}
	_, err := c.catalog.Deregister(dereg, &consulapi.WriteOptions{})
	if err != nil {
		return fmt.Errorf("failed to deregister consul service: %w", err)
	}
	return nil
}
