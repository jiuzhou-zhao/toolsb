package main

import (
	"context"

	"github.com/jiuzhou-zhao/go-fundamental/loge"
	"github.com/jiuzhou-zhao/go-fundamental/servicetoolset"
	"github.com/jiuzhou-zhao/go-fundamental/tracing"
	"github.com/sbasestarter/proto-repo/gen/protorepo-tool-go"
	"github.com/sbasestarter/toolsb/internal/config"
	"github.com/sbasestarter/toolsb/internal/server"
	"github.com/sgostarter/libconfig"
	"github.com/sgostarter/liblog"
	"google.golang.org/grpc"
)

func main() {
	logger, err := liblog.NewZapLogger()
	if err != nil {
		panic(err)
	}
	loggerChain := loge.NewLoggerChain()
	loggerChain.AppendLogger(tracing.NewTracingLogger())
	loggerChain.AppendLogger(logger)
	loge.SetGlobalLogger(loge.NewLogger(loggerChain))

	var cfg config.Config
	_, err = libconfig.Load("config", &cfg)
	if err != nil {
		loge.Fatalf(context.Background(), "load config failed: %v", err)
		return
	}

	serverImpl := server.NewServer(&cfg)
	serviceToolset := servicetoolset.NewServerToolset(context.Background(), loggerChain)
	_ = serviceToolset.CreateGRpcServer(&cfg.GRpcServerConfig, nil, func(s *grpc.Server) {
		toolpb.RegisterUserServiceServer(s, serverImpl)
	})
	cfg.HttpServerConfig.Handler = serverImpl.GetHTTPHandler()
	_ = serviceToolset.CreateHttpServer(&cfg.HttpServerConfig)
	serviceToolset.Wait()
}
