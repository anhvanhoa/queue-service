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
	Addr        string
	DB          int
	Password    string
	MaxIdle     int
	MaxActive   int
	IdleTimeout int
	Network     string
	Concurrency int
	Queues      map[string]int
}

type GrpcClient struct {
	ServerAddress string
	Timeout       time.Duration
	MaxRetries    int
	KeepAlive     *keepalive.ClientParameters
}

type Env struct {
	MODE_ENV     string
	NAME_APP     string
	QUEUE        *queue
	GRPC_CLIENTS map[string]*GrpcClient
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
	return strings.ToLower(env.MODE_ENV) == "production"
}
