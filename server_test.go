package microBase

import (
	"context"
	"fmt"
	go_micro_srv_example "github.com/Gitforxuyang/microBase/proto/example"
	"github.com/go-errors/errors"
	"os"
	"testing"
	"time"
)

type Example struct{}

// Call is a single request handler called via client.Call or the generated client code
func (e *Example) Call(ctx context.Context, req *go_micro_srv_example.Request, rsp *go_micro_srv_example.Response) error {
	rsp.Msg = "Hello " + req.Name
	//return nil
	return errors.New("123")
}

// Stream is a server side stream handler called via client.Stream or the generated client code
func (e *Example) Stream(ctx context.Context, req *go_micro_srv_example.StreamingRequest, stream go_micro_srv_example.Example_StreamStream) error {

	for i := 0; i < int(req.Count); i++ {
		if err := stream.Send(&go_micro_srv_example.StreamingResponse{
			Count: int64(i),
		}); err != nil {
			return err
		}
	}

	return nil
}

// PingPong is a bidirectional stream handler called via client.Stream or the generated client code
func (e *Example) PingPong(ctx context.Context, stream go_micro_srv_example.Example_PingPongStream) error {
	for {
		req, err := stream.Recv()
		if err != nil {
			return err
		}
		if err := stream.Send(&go_micro_srv_example.Pong{Stroke: req.Stroke}); err != nil {
			return err
		}
	}
}

var service MicroService

func TestMain(m *testing.M) {
	service = MicroInit()
	exitCode := m.Run()
	os.Exit(exitCode)
}

func TestMicroInit(t *testing.T) {
	//go_micro_srv_example.RegisterExampleHandler(service.Server(), new(Example))
	//service.Run()
	client := service.Client()
	srv := go_micro_srv_example.NewExampleService("127.0.0.1:7001", client)
	rsp, err := srv.Call(context.TODO(), &go_micro_srv_example.Request{Name: "123123"})
	if err != nil {
		fmt.Printf(err.Error())
	}
	fmt.Println(rsp)
	time.Sleep(time.Second * 2)
}

func TestClient(t *testing.T) {
	//fmt.Println(1111)
}
