package services

import (
	"errors"
	"fmt"
	"time"

	"github.com/Strum355/log"
	"github.com/hashicorp/consul/api"
	"github.com/spf13/viper"
)

type ConsulService struct {
	ID          string
	ConsulHost  string
	ConsulToken string
	ServiceAddr string
	Port        int
	TTL         time.Duration
	client      *api.Client
}

func (c *ConsulService) Setup() error {
	config := api.Config{
		Address: c.ConsulHost,
		Token:   c.ConsulToken,
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
		ID:      c.ID,
		Name:    "fail2rest",
		Address: c.ServiceAddr,
		Port:    c.Port,
		Check: &api.AgentServiceCheck{
			DeregisterCriticalServiceAfter: (time.Second * 10).String(),
			TTL:                            c.TTL.String(),
		},
	}

	err := c.getSharedSecret()
	if err != nil {
		return err
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
			if viper.GetString("fail2rest.secret") == "" {
				health = api.HealthCritical
			}
			err := c.client.Agent().UpdateTTL("service:"+c.ID, "", health)
			if err != nil {
				log.WithError(err).Error("failed to update TTL")
			}
		}
	}()
}

func (c *ConsulService) getSharedSecret() error {
	fn := func() error {
		path := viper.GetString("consul.fail2rest")
		kv, _, err := c.client.KV().Get(path, &api.QueryOptions{})
		if err != nil {
			return err
		}

		if kv == nil {
			return errors.New(fmt.Sprintf("key %s not set", path))
		}

		viper.Set("fail2rest.secret", string(kv.Value))
		return nil
	}

	count := 4
	var err error
	for ; count > 0; count-- {
		err = fn()
		if err == nil {
			return nil
		}
		log.WithFields(log.Fields{
			"limit": 4,
			"count": count,
		}).WithError(err).Error("failed to get shared secret")
		time.Sleep(time.Second * 3)
	}
	return err
}
