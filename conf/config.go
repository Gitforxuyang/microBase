package conf

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
