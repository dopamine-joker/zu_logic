package main

import (
	"context"
	"fmt"
	"log"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-basic/uuid"

	"github.com/dopamine-joker/zu_logic/handle/handle"
	"github.com/dopamine-joker/zu_logic/misc"
)

func process() {
	var err error
	serverId := fmt.Sprintf("logic-%s", uuid.New())
	if err = handle.InitRpcServer(serverId); err != nil {
		panic(err)
	}
}

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	defer stop()
	misc.Init()
	process()
	handle.TaskAddOrder(ctx)
	handle.TaskUpdateOrder(ctx)
	<-ctx.Done()
	handle.StopServer()
	misc.Logger.Info("zu_logic exit")
	log.Println("zu_logic exit")
	// 等待资源回收
	time.Sleep(2 * time.Second)
}
