package lib

import (
	"context"
	"fmt"
	clientv3 "go.etcd.io/etcd/client/v3"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"time"
)

type EtcdClient struct {
	client *clientv3.Client
	Services []*ServiceInfo
}

type ServiceInfo struct {
	ServiceID string
	ServiceName string
	ServiceAddr string
}

func NewClient() *EtcdClient {
	config := clientv3.Config{
		Endpoints: []string{"localhost:2379", "localhost:22379", "localhost:32379"},
		DialTimeout: 10 * time.Second,
	}
	client, err := clientv3.New(config)
	if err != nil {
		log.Fatal("etcd客户端创建失败！")
	}

	return &EtcdClient{
		client: client,
	}

}

func (s *EtcdClient) LoadService() error {
	kv := clientv3.NewKV(s.client)
	res, err := kv.Get(context.Background(), "/service", clientv3.WithPrefix())
	if err != nil {
		log.Fatal(err)
	}
	for _, v := range res.Kvs {
		fmt.Println(string(v.Key))
		s.parseService(v.Key, v.Value)
	}

	return err

}

func (s *EtcdClient) parseService(key []byte, value []byte) {

	res := regexp.MustCompile("/service/(.+)/(\\w+)")
	if res.Match(key) {
		idAndName := res.FindSubmatch(key)
		serviceId := idAndName[1]
		serviceName := idAndName[2]
		s.Services = append(s.Services, &ServiceInfo{ServiceID: string(serviceId), ServiceName: string(serviceName), ServiceAddr: string(value)})

	}

}

func (s *EtcdClient) GetService(serviceName string, method string, encodeFunc EncodeRequestFunc) Endpoint {
	
	//for _, service := range s.Services {
	//	if serviceName == service.ServiceName {
	//		return func(ctx context.Context, requestParam interface{}) (respenseResult interface{}, err error) {
	//			httpClient := http.DefaultClient
	//			httpRequest, err := http.NewRequest(method, "http://" + service.ServiceAddr, nil)
	//			err = encodeFunc(ctx, httpRequest, requestParam)
	//			if err != nil {
	//				return nil, err
	//			}
	//
	//			resq, err := httpClient.Do(httpRequest)
	//			defer resq.Body.Close()
	//			if err != nil {
	//				return nil, err
	//			}
	//			body, err := ioutil.ReadAll(resq.Body)
	//			if err != nil {
	//				return nil, err
	//			}
	//
	//			return string(body), nil
	//		}
	//	}
	//}

	lb := NewloadBalance(s.Services)
	serviceInstance := lb.getByRand(serviceName)
	if serviceInstance == nil {
		return nil

	}


	return func(ctx context.Context, requestParam interface{}) (respenseResult interface{}, err error) {
		httpClient := http.DefaultClient
		httpRequest, err := http.NewRequest(method, "http://" + serviceInstance.ServiceAddr, nil)
		err = encodeFunc(ctx, httpRequest, requestParam)
		if err != nil {
			return nil, err
		}

		resq, err := httpClient.Do(httpRequest)
		defer resq.Body.Close()
		if err != nil {
			return nil, err
		}
		body, err := ioutil.ReadAll(resq.Body)
		if err != nil {
			return nil, err
		}

		return string(body), nil
	}

}