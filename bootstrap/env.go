package bootstrap

import (
	"strings"

	"github.com/anhvanhoa/service-core/bootstrap/config"
	"github.com/anhvanhoa/service-core/domain/grpc_client"
)

type queue struct {
	Addr        string         `mapstructure:"addr"`
	Db          int            `mapstructure:"db"`
	Password    string         `mapstructure:"password"`
	MaxIdle     int            `mapstructure:"max_idle"`
	MaxActive   int            `mapstructure:"max_active"`
	IdleTimeout int            `mapstructure:"idle_timeout"`
	Network     string         `mapstructure:"network"`
	Concurrency int            `mapstructure:"concurrency"`
	Queues      map[string]int `mapstructure:"queues"`
}

type Env struct {
	NodeEnv         string                    `mapstructure:"node_env"`
	NameService     string                    `mapstructure:"name_service"`
	HostGrpc        string                    `mapstructure:"host_grpc"`
	PortGrpc        int                       `mapstructure:"port_grpc"`
	IntervalCheck   string                    `mapstructure:"interval_check"`
	TimeoutCheck    string                    `mapstructure:"timeout_check"`
	Queue           *queue                    `mapstructure:"queue"`
	MailServiceAddr string                    `mapstructure:"mail_service_addr"`
	GrpcClients     []*grpc_client.ConfigGrpc `mapstructure:"grpc_clients"`
}

func NewEnv(env any) {
	setting := config.DefaultSettingsConfig()
	if setting.IsProduction() {
		setting.SetPath("/config")
		setting.SetFile("queue_service.config")
	} else {
		setting.SetFile("dev.config")
	}
	config.NewConfig(setting, env)
}

func (env *Env) IsProduction() bool {
	return strings.ToLower(env.NodeEnv) == "production"
}
