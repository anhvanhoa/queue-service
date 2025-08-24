package grpc_client

import (
	"fmt"
	"sync"
)

type ClientFactory struct {
	clients map[string]*Client
	mutex   sync.RWMutex
}

func NewClientFactory(config ...*Config) *ClientFactory {
	cf := &ClientFactory{
		clients: make(map[string]*Client),
	}
	if len(config) > 0 {
		for _, c := range config {
			cf.CreateClient(c)
		}
	}
	return cf
}

func (cf *ClientFactory) CreateClient(config *Config) (*Client, error) {
	cf.mutex.Lock()
	defer cf.mutex.Unlock()

	if client, exists := cf.clients[config.Name]; exists {
		if client.IsConnected() {
			return client, nil
		}
		client.Close()
	}

	client, err := NewClient(config)
	if err != nil {
		return nil, fmt.Errorf("không thể tạo client '%s': %w", config.Name, err)
	}
	if client.IsConnected() {
		fmt.Printf("Kết nối thành công: %s - Host: %s\n", config.Name, client.conn.Target())
	} else {
		fmt.Printf("Kết nối thất bại: %s - Host: %s\n", config.Name, client.conn.Target())
	}
	cf.clients[config.Name] = client
	return client, nil
}

func (cf *ClientFactory) GetClient(name string) *Client {
	cf.mutex.RLock()
	defer cf.mutex.RUnlock()
	return cf.clients[name]
}

func (cf *ClientFactory) CloseAll() error {
	cf.mutex.Lock()
	defer cf.mutex.Unlock()

	var lastErr error
	for name, client := range cf.clients {
		if err := client.Close(); err != nil {
			lastErr = fmt.Errorf("không thể đóng client '%s': %w", name, err)
		}
		delete(cf.clients, name)
	}

	return lastErr
}

func (cf *ClientFactory) CloseClient(name string) error {
	cf.mutex.Lock()
	defer cf.mutex.Unlock()

	if client, exists := cf.clients[name]; exists {
		err := client.Close()
		delete(cf.clients, name)
		return err
	}

	return fmt.Errorf("không tìm thấy client '%s'", name)
}

func (cf *ClientFactory) ListClients() []string {
	cf.mutex.RLock()
	defer cf.mutex.RUnlock()

	var names []string
	for name := range cf.clients {
		names = append(names, name)
	}
	return names
}

func (cf *ClientFactory) HealthCheck() map[string]bool {
	cf.mutex.RLock()
	defer cf.mutex.RUnlock()

	health := make(map[string]bool)
	for name, client := range cf.clients {
		health[name] = client.IsConnected()
	}
	return health
}
