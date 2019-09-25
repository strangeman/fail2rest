package services

import (
	"fmt"
	"time"

	"github.com/hashicorp/consul/api"
	"github.com/Strum355/log"
)

type ConsulService struct {
	ID string
	ConsulHost string
	ConsulToken string
	ServiceAddr string
	Port int
	TTL time.Duration
	Secret string
	client *api.Client
}

func (c *ConsulService) Setup() error {
	config := api.Config{
		Address: c.ConsulHost,
		Token: c.ConsulToken,
	}
	client, err := api.NewClient(&config)
	if err != nil {
		return err
	}

	c.client = client
	return nil
}

func (c *ConsulService) Register() error {
	c.ID = fmt.Sprintf("fail2rest@%s", c.ServiceAddr)

	service := &api.AgentServiceRegistration{
		ID: c.ID,
		Name: "fail2rest",
		Address: c.ServiceAddr,
		Port: c.Port,
		Check: &api.AgentServiceCheck{
			DeregisterCriticalServiceAfter: (time.Second * 10).String(),
			TTL:                            c.TTL.String(),
		},
	}

	registered := c.client.Agent().ServiceRegister(service)

	c.updateTTL()

	return registered
}

func (c *ConsulService) updateTTL() {
	go func() {
		ticker := time.NewTicker(c.TTL / 2)
		for range ticker.C {
			health := api.HealthPassing
			if c.Secret == "" {
				health = api.HealthCritical
			}
			err := c.client.Agent().UpdateTTL("service:"+c.ID, "", health)
			if err != nil {
				log.WithError(err).Error("failed to update TTL")
			}
		}
	}()
}
