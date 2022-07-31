package main

import (
	_ "bookstore/internal/store"
	"bookstore/server"
	"bookstore/store/factory"
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {

	store, err := factory.New("mem")
	if err != nil {
		panic(err)
	}

	// 创建http服务实例
	storeServer := server.NewBookStoreServer(":8080", store)
	// 运行http服务
	errChan, err := storeServer.ListenAndServe()
	if err != nil {
		log.Println("web server start failed: ", err)
		return
	}
	log.Println("web server start success")

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)

	// 监听来自errChan以及c的事件
	select {
	case err = <-errChan:
		log.Println("web server run failed: ", err)
		return
	case <-c:
		log.Println("bookstore program is exiting...")
		ctx, cancelFunc := context.WithTimeout(context.Background(), time.Second)
		defer cancelFunc()
		err = storeServer.Shutdown(ctx)
	}

	if err != nil {
		log.Println("bookstore program exit error: ", err)
		return
	}
	log.Println("bookstore program exit success")

}
