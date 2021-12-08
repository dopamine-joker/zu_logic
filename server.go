package main

import (
	"context"
	"log"
	"os/signal"
	"syscall"
	"time"
	"zu_logic/misc"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	defer stop()
	misc.Init()
	<-ctx.Done()
	misc.Logger.Info("zu_logic exit")
	log.Println("zu_logic exit")
	// 等待资源回收
	time.Sleep(2 * time.Second)
}
