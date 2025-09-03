package bootstrap

import (
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/viper"
	"google.golang.org/grpc/keepalive"
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

type GrpcClient struct {
	ServerAddress string                      `mapstructure:"server_address"`
	Timeout       time.Duration               `mapstructure:"timeout"`
	MaxRetries    int                         `mapstructure:"max_retries"`
	KeepAlive     *keepalive.ClientParameters `mapstructure:"keep_alive"`
}

type Env struct {
	NodeEnv     string                 `mapstructure:"node_env"`
	NameServer  string                 `mapstructure:"name_server"`
	Queue       *queue                 `mapstructure:"queue"`
	GrpcClients map[string]*GrpcClient `mapstructure:"grpc_clients"`
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
