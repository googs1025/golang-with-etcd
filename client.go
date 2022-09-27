package main

import (
	"context"
	"etcd-practice/lib"
	"etcd-practice/service"
	"fmt"
	"log"
)

func main() {
	client := lib.NewClient()
	err := client.LoadService()
	if err != nil {
		log.Fatal(err)
	}
	//fmt.Println(client.Services)

	//for _,v := range client.Services {
	//	fmt.Printf("serviceId:%s, serviceName:%s, serviceAddr:%s", v.ServiceID, v.ServiceName, v.ServiceAddr)
	//}

	endpoint := client.GetService("productservice", "GET", service.ProdEncodeFunc)
	resp, err := endpoint(context.Background(), service.ProductRequest{ProductId: 106})
	if err !=nil {
		log.Fatal(err)
	}
	fmt.Println(resp)

}
