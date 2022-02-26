package handle

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"log"
	"net"

	pb "github.com/dopamine-joker/zu_logic/proto"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"

	"github.com/dopamine-joker/zu_logic/misc"
)

var stopCh chan struct{}

func StopServer() {
	close(stopCh)
}

func InitRpcServer(serverId string) error {
	stopCh = make(chan struct{})
	localIP := misc.GetLocalIP()
	// 获取本机ip,用于在etcd注册服务
	misc.Logger.Info("local IP", zap.String("IP", localIP))
	for _, port := range misc.Conf.Logic.RpcPort {
		addr := fmt.Sprintf("%s:%d", localIP, port)
		go createGrpcServer(misc.Network, addr, serverId)
		misc.Logger.Info("logic start run", zap.String("addr", fmt.Sprintf("%s:%d", localIP, port)))
		log.Printf("logic start run at --> %s", fmt.Sprintf("%s:%d", localIP, port))
	}
	return nil
}

//createGrpcServer 创建rpc server
func createGrpcServer(network, addr string, serverId string) {
	grpcServer := grpc.NewServer(
		grpc.MaxRecvMsgSize(8*1024*1024),
		grpc.MaxSendMsgSize(8*1024*1024),
		grpc.UnaryInterceptor(otelgrpc.UnaryServerInterceptor()),
		grpc.StreamInterceptor(otelgrpc.StreamServerInterceptor()),
	)
	rpcServer := NewRpcLogicServer(addr)
	pb.RegisterRpcLogicServiceServer(grpcServer, rpcServer)
	listener, err := net.Listen(network, rpcServer.Addr)
	if err != nil {
		misc.Logger.Error("net.Listen fail when register server", zap.Error(err))
		panic(err)
	}
	etcdRegister, err := NewRegister(misc.Conf.EtcdCfg.Host, misc.Conf.EtcdCfg.BasePath, misc.Conf.EtcdCfg.ServerPathLogic, 5)
	if err != nil {
		misc.Logger.Error("NewRegister err", zap.Error(err))
		panic(err)
	}
	if err = etcdRegister.Register(context.Background(), rpcServer, 5); err != nil {
		misc.Logger.Error("Register to etcd fail", zap.Error(err))
		panic(err)
	}
	go func() {
		<-stopCh
		misc.Logger.Info("grpcServer stop", zap.String("host", addr))
		grpcServer.Stop()
		etcdRegister.StopServe()
	}()
	if err = grpcServer.Serve(listener); err != nil {
		misc.Logger.Error("grpcServer serve failed")
		panic(err)
	}
}
