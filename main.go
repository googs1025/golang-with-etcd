package main

import (
	"context"
	"etcd-practice/lib"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"net/http"
)

// 服务注册：需要把ip地址与服务名称写进etcd。

func main()  {

	router := mux.NewRouter()

	router.HandleFunc("/product/{id:\\d+}", func(writer http.ResponseWriter, request *http.Request) {
		vars := mux.Vars(request)
		str := "get product ByID:" + vars["id"]
		writer.Write([]byte(str))
	})

	s := lib.NewEtcdClientService()
	serviceID := "p3"
	serviceName := "productservice"
	serviceAddr := "127.0.0.1"
	servicePort := 8081

	httpServer := &http.Server{
		Addr: ":" + strconv.Itoa(servicePort),
		Handler: router,
	}

	errC := make(chan error, 1)

	// 服务使用另一个goroutine执行，主goroutine后面用来优雅退出。
	go func() {
		// 服务注册
		err := s.RegService(serviceID, serviceName, serviceAddr + ":" + strconv.Itoa(servicePort))
		if err != nil {
			errC <-err
			return
		}
		// 启动server
		err = httpServer.ListenAndServe()
		if err != nil {
			errC <-err
			return
		}
	}()

	// 监听退出chan
	go func() {
		stopC := make(chan os.Signal, 1)
		signal.Notify(stopC,syscall.SIGINT,syscall.SIGTERM)
		errC <-fmt.Errorf("%s", <-stopC)
	}()

	// 优雅退出流程

	getErr :=<-errC
	fmt.Println("需要从etcd中注销服务")
	err := s.UnRegService(serviceID)
	if err != nil {
		fmt.Println(err)
	}

	err = httpServer.Shutdown(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	log.Println("发生异常，服务正在停止。。。")
	log.Fatal(getErr)





}
