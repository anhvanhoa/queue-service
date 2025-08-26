package discovery

import (
	"fmt"
	"net"
	"queue-service/bootstrap"
	loggerI "queue-service/domain/service/logger"

	"github.com/hashicorp/consul/api"
)

type Discovery interface {
	Register(serviceName string, servicePort int) error
	GetService(serviceName string) (string, error)
	Close(id string) error
}

type DiscoveryConfig struct {
	ServiceName    string
	ServiceAddress string
	ServicePort    int
	ServiceHost    string
	ServiceID      string
	ServiceTags    []string
	ServiceMeta    map[string]string
}

type DiscoveryImpl struct {
	client *api.Client
	log    loggerI.Log
	env    *bootstrap.Env
}

func NewDiscovery(log loggerI.Log, env *bootstrap.Env) (*DiscoveryImpl, error) {
	client, err := api.NewClient(api.DefaultConfig())
	if err != nil {
		return nil, err
	}
	return &DiscoveryImpl{
		client: client,
		log:    log,
		env:    env,
	}, nil
}

func (d *DiscoveryImpl) Register(serviceName string) error {
	registration := new(api.AgentServiceRegistration)
	addrs, _ := net.InterfaceAddrs()
	serviceAddress := d.env.HOST_GRPC
	for _, addr := range addrs {
		if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() && ipnet.IP.To4() != nil && d.env.IsProduction() {
			serviceAddress = ipnet.IP.String()
		}
	}
	registration.Name = serviceName
	registration.Port = d.env.PORT_GRPC
	registration.Address = serviceAddress
	registration.Check = &api.AgentServiceCheck{
		GRPC:     fmt.Sprintf("%s:%d", serviceAddress, d.env.PORT_GRPC),
		Interval: d.env.INTERVAL_CHECK,
		Timeout:  d.env.TIMEOUT_CHECK,
	}

	if err := d.client.Agent().ServiceRegister(registration); err != nil {
		d.log.Error(fmt.Sprintf("Failed to register service: %v", err))
		return err
	}
	d.log.Info(fmt.Sprintf("Service %s registered with address: %s:%d", serviceName, serviceAddress, d.env.PORT_GRPC))
	return nil
}

func (d *DiscoveryImpl) GetService(serviceName string) (*api.ServiceEntry, error) {
	services, _, err := d.client.Health().Service(serviceName, "", true, nil)
	if err != nil {
		d.log.Error("Failed to get service", err)
		return nil, err
	}
	if len(services) == 0 {
		return nil, fmt.Errorf("service %s not found", serviceName)
	}
	service := services[0]
	return service, nil
}

func (d *DiscoveryImpl) Close(id string) error {
	return d.client.Agent().ServiceDeregister(id)
}
