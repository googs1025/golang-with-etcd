package main

import (
	"crypto/md5"
	"encoding/hex"
	"flag"
	"fmt"
	"github.com/jpillora/overseer"
	"gopkg.in/ini.v1"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"
)


// 使用server 监听配置更新情况。本次使用ini文件

func main()  {
	port:=flag.Int("p",0,"服务端口")
	flag.Parse()
	if *port==0{
		log.Fatal("请指定端口")
	}
	cfg,err:=ini.Load("my.ini")
	if err!=nil{
		log.Fatal(err)
	}
	mux:=http.NewServeMux()

	// 显示配置文件
	mux.HandleFunc("/",func(writer http.ResponseWriter, request *http.Request) {
		dbUser:=cfg.Section("db").Key("db_user").Value()
		dbPass:=cfg.Section("db").Key("db_pass").Value()
		_, _ = writer.Write([]byte("<h1>" + dbUser + "</h1>"))
		_, _ = writer.Write([]byte("<h1>" + dbPass + "</h1>"))
	})
	// 重新加载ini文件
	mux.HandleFunc("/reload",func(writer http.ResponseWriter, request *http.Request) {
		newCfg,_:=ini.Load("my.ini")
		cfg=newCfg
	})


	// 建立server
	server:=&http.Server{
		Addr:":"+strconv.Itoa(*port),
		Handler:mux,
	}
	// 重启server
	prog:= func(state overseer.State) {
		_ = server.Serve(state.Listener)
	}

	// 平滑重启
	errChan:=make(chan error)
	go func() {
		overseer.Run(overseer.Config{
			Program: prog,
			TerminateTimeout:time.Second*2,
			Address: ":"+strconv.Itoa(*port),
		})
	}()

	//监听信号
	go func() {
		sigC := make(chan os.Signal)
		signal.Notify(sigC,syscall.SIGINT,syscall.SIGTERM)
		errChan<-fmt.Errorf("%s",<-sigC)
	}()

	//监控配置文件变化
	go func() {
		fileMd5, err := getFileMD5("my.ini")
		if err != nil {
			log.Println(err)
			return
		}
		for{
			newMd5, err := getFileMD5("my.ini")
			if err != nil {
				log.Println(err)
				break
			}
			if strings.Compare(newMd5,fileMd5) != 0 {
				fileMd5 = newMd5
				fmt.Println("文件发生了变化")
				// 如果文件发生变化，server执行重启操作
				overseer.Restart()
			}
			time.Sleep(time.Second*2)
		}
	}()

	getErr := <-errChan

	log.Println(getErr)
}

// 利用md5值判断文件是否修改
func getFileMD5(filePath string) (string,error)  {
	file, err := os.Open(filePath)
	hash := md5.New()
	if _, err = io.Copy(hash, file); err != nil {
		return "",err
	}
	hashInBytes := hash.Sum(nil)[:16]
	return hex.EncodeToString(hashInBytes),nil
}

