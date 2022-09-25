package main

import (
	"etcd-practice/lib"
	"fmt"
	"log"
)

func main() {
	client := lib.NewClient()
	err := client.GetService()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(client.Services)

	for _,v := range client.Services {
		fmt.Printf("serviceId:%s, serviceName:%s, serviceAddr:%s", v.ServiceID, v.ServiceName, v.ServiceAddr)
	}

}
