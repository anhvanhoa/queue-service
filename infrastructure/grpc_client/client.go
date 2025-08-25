package grpc_client

import (
	"fmt"
	loggerI "queue-service/domain/service/logger"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"
)

type Client struct {
	conn   *grpc.ClientConn
	config *Config
	log    loggerI.Log
}

type Config struct {
	Name          string
	ServerAddress string
	Timeout       time.Duration
	MaxRetries    int
	KeepAlive     *keepalive.ClientParameters
}

func NewClient(config *Config, log loggerI.Log) (*Client, error) {
	if config == nil {
		config = &Config{
			ServerAddress: "localhost:50051",
			Timeout:       30 * time.Second,
			MaxRetries:    3,
			KeepAlive: &keepalive.ClientParameters{
				Time:                10 * time.Second,
				Timeout:             20 * time.Second,
				PermitWithoutStream: true,
			},
		}
	}

	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithKeepaliveParams(*config.KeepAlive),
	}

	conn, err := grpc.NewClient(config.ServerAddress, opts...)
	if err != nil {
		return nil, fmt.Errorf("không thể kết nối đến máy chủ gRPC: %w", err)
	}

	return &Client{
		conn:   conn,
		config: config,
		log:    log,
	}, nil
}

func (c *Client) GetConnection() *grpc.ClientConn {
	return c.conn
}

func (c *Client) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}

func (c *Client) IsConnected() bool {
	c.log.Info(fmt.Sprintf("Trạng thái kết nối: %s", c.conn.GetState().String()))
	return c.conn != nil
}
