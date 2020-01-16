package microBase

import (
	"fmt"
	"github.com/Gitforxuyang/microBase/conf"
	"github.com/Gitforxuyang/microBase/trace"
	"github.com/Gitforxuyang/microBase/util"
	"github.com/Gitforxuyang/microBase/wrapper"
	"github.com/micro/go-micro"
	memory2 "github.com/micro/go-micro/broker/memory"
	"github.com/micro/go-micro/client"
	"github.com/micro/go-micro/client/selector/static"
	"github.com/micro/go-micro/registry/memory"
	"github.com/micro/go-micro/server"
	"github.com/micro/go-micro/service/grpc"
	"github.com/micro/go-micro/util/log"
)

type MicroService interface {
	Server() server.Server
	Run()
	Client() client.Client
}

type microService struct {
	s micro.Service
}

func (m *microService) Server() server.Server {
	return m.s.Server()
}

func (m *microService) Run() {
	if err := m.s.Run(); err != nil {
		log.Fatal(err)
		panic("server run err")
	}
}

func (m *microService) Client() client.Client {
	return m.s.Client()
}

var (
	BaseConfig conf.Config
)

func MicroInit() MicroService {
	util.InitLog()
	BaseConfig = conf.InitConfig()
	//= conf.GetConfig(env)
	//err := config.Load(
	//	file.NewSource(
	//		file.WithPath("./conf/config.default.json"),
	//	),
	//	file.NewSource(
	//		file.WithPath(fmt.Sprintf("./conf/config.%s.json", env)),
	//	),
	//)
	//Must(err)
	//val := config.Get("port")
	//port := val.Int(7001)
	//name := config.Get("name").String("server")
	//version := config.Get("version").String("0.0.1")
	util.ErrInit(BaseConfig.ServerConfig.Port)
	tracer, closer, err := trace.NewTracer(
		fmt.Sprintf("%s_%s", BaseConfig.ServerConfig.ServerName, BaseConfig.ServerConfig.Env),
		BaseConfig.Traceing.Endpoint)
	util.Must(err)
	// New Service
	service := grpc.NewService(
		micro.Name(BaseConfig.ServerConfig.ServerName),
		micro.Version(BaseConfig.ServerConfig.Version),
		micro.Address(fmt.Sprintf("0.0.0.0:%d", BaseConfig.ServerConfig.Port)),
		micro.Registry(memory.NewRegistry()),
		micro.Broker(memory2.NewBroker()),
		micro.Selector(static.NewSelector()),
		//micro.Flags(cli.StringFlag{
		//	Name:   "ENV",
		//	EnvVar: "ENV",
		//	Value:  "local",
		//}),
		micro.WrapHandler(
			wrapper.NewTraceWrapper(tracer),
			wrapper.NewLogWrapper()),
		micro.WrapClient(wrapper.NewClientWrapper()),
		micro.WrapCall(wrapper.NewCallTraceWrapper(tracer)),
		micro.BeforeStart(func() error {
			return nil
		}),
		micro.AfterStart(func() error {
			log.Infof("server started listen in :%d", BaseConfig.ServerConfig.Port)
			return nil
		}),
		micro.BeforeStop(func() error {
			//log.Info("before end")
			closer.Close()
			return nil
		}),
		micro.AfterStop(func() error {
			log.Info("server closed")
			return nil
		}),
	)
	//service.Init(micro.Action(func(cli *cli.Context) {
	//
	//}))
	return &microService{s: service}
}
