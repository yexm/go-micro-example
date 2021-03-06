package main

import (
	"time"

	"github.com/micro/go-micro/service/grpc"
	"github.com/micro/go-micro/util/log"
	"github.com/micro/go-micro/web"
	"github.com/micro/go-plugins/registry/consul"
	"github.com/opentracing/opentracing-go"

	"go-micro-example/api/api/router"
	"go-micro-example/pkg/tracer"
)

const (
	Version = "v0.1.0"
	Name    = "com.example.api.api"
)

func main() {
	reg := consul.NewRegistry()

	t, io, err := tracer.NewTracer(Name, "")
	if err != nil {
		log.Fatal(err)
	}
	defer io.Close()
	opentracing.SetGlobalTracer(t)

	// create new web service
	service := web.NewService(
		web.Name(Name),
		web.Version(Version),
		web.MicroService(grpc.NewService()),
		web.RegisterInterval(time.Duration(5)),
		web.RegisterTTL(time.Duration(10)),
		web.Registry(reg),
	)

	// initialise service
	if err := service.Init(); err != nil {
		log.Fatal(err)
	}

	// init router
	r := router.InitRouter(service.Options().Service.Client())

	service.Handle("/", r)
	// run service
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
