package conf

import (
	"flag"
	"fmt"
	"github.com/Gitforxuyang/microBase/util"
	"github.com/micro/go-micro/config"
	"github.com/micro/go-micro/config/source/file"
	"path/filepath"
)

//链路追踪配置
type traceingConfig struct {
	Endpoint string
}

//服务基础配置
type serverConfig struct {
	Port       int
	ServerName string
	Env        string
	Version    string
}

type Config struct {
	Traceing     traceingConfig
	ServerConfig serverConfig
}

func GetConfig(env string) Config {
	switch env {
	case "local":
		return initLocalConfig()
	default:
		panic("err")
	}
}

func initLocalConfig() Config {
	config := Config{}
	t := traceingConfig{
		Endpoint: "http://127.0.0.1:14268/api/traces",
	}
	config.Traceing = t
	return config
}

func InitConfig() Config {
	env := *flag.String("ENV", "local", "环境变量")
	if env != "prod" && env != "dev" && env != "local" {
		panic("ENV只能是local dev prod之一")
	}
	sp := string(filepath.Separator)
	appPath, _ := filepath.Abs(filepath.Dir(filepath.Join("."+sp, sp)))
	pt := filepath.Join(appPath, "/conf")
	err := config.Load(
		file.NewSource(
			file.WithPath(pt+"/config.default.json"),
		),
		file.NewSource(
			file.WithPath(fmt.Sprintf(pt+"/config.%s.json", env)),
		),
	)
	util.Must(err)
	cfg := GetConfig(env)
	cfg.ServerConfig.Port = config.Get("port").Int(7001)
	cfg.ServerConfig.Version = config.Get("port").String("0.0.1")
	cfg.ServerConfig.ServerName = config.Get("name").String("server")
	return cfg
}
