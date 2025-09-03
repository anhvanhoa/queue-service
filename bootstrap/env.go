package bootstrap

import (
	"log"
	"os"
	"path/filepath"
	"queue-service/infrastructure/grpc_client"
	"strings"

	"github.com/spf13/viper"
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
	NodeEnv         string                `mapstructure:"node_env"`
	NameService     string                `mapstructure:"name_service"`
	HostGrpc        string                `mapstructure:"host_grpc"`
	PortGrpc        int                   `mapstructure:"port_grpc"`
	IntervalCheck   string                `mapstructure:"interval_check"`
	TimeoutCheck    string                `mapstructure:"timeout_check"`
	Queue           *queue                `mapstructure:"queue"`
	MailServiceAddr string                `mapstructure:"mail_service_addr"`
	GrpcClients     []*grpc_client.Config `mapstructure:"grpc_clients"`
}

func NewEnv(env *Env) {
	absPath, err := filepath.Abs("./")
	if err != nil {
		log.Fatal("Lỗi khi lấy đường dẫn tuyệt đối:", err)
	}

	mode := os.Getenv("ENV_MODE")
	viper.SetConfigType("yaml")
	if mode == "production" {
		viper.SetConfigName("prod.config")
	} else {
		viper.SetConfigName("dev.config")
	}
	viper.AddConfigPath(absPath)
	err = viper.ReadInConfig()
	if err != nil {
		panic("Lỗi khi đọc file cấu hình, " + err.Error())
	}

	err = viper.UnmarshalExact(env)
	if err != nil {
		panic("Lỗi khi phân tích file cấu hình, " + err.Error())
	}
}

func (env *Env) IsProduction() bool {
	return strings.ToLower(env.NodeEnv) == "production"
}
