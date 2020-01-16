package conf

import (
	"fmt"
	"github.com/Gitforxuyang/microBase"
	"github.com/micro/go-micro/config"
	"github.com/micro/go-micro/config/source/file"
	"path/filepath"
)

type traceingConfig struct {
	Endpoint string
}

type Config struct {
	Traceing traceingConfig
}

var (
	local Config = initLocalConfig()
)

func GetConfig(env string) Config {
	switch env {
	case "local":
		return local
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

func InitFileConfig(env string) {
	sp := string(filepath.Separator)
	appPath, _ := filepath.Abs(filepath.Dir(filepath.Join("."+sp, sp)))
	pt := filepath.Join(appPath, "conf")
	err := config.Load(
		file.NewSource(
			file.WithPath(pt+"config.default.json"),
		),
		file.NewSource(
			file.WithPath(fmt.Sprintf(pt+"config.%s.json", env)),
		),
	)
	microBase.Must(err)
}
