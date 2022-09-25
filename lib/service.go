package lib

import (
	"context"
	"fmt"
	clientv3 "go.etcd.io/etcd/client/v3"
	"log"
	"time"
)

// EtcdClientService etcd客户端对象
type EtcdClientService struct {
	client *clientv3.Client
}

// NewEtcdClientService 构造函数
func NewEtcdClientService() *EtcdClientService {
	config := clientv3.Config{
		Endpoints: []string{"localhost:2379", "localhost:22379", "localhost:32379"},
		DialTimeout: 10 * time.Second,
	}
	client, err := clientv3.New(config)
	if err != nil {
		log.Fatal("etcd客户端创建失败！")
	}

	return &EtcdClientService{
		client: client,
	}
}


func (s *EtcdClientService) RegService(id string, name string, addr string) error {
	kv := clientv3.NewKV(s.client)
	key_prefix := "/service/"
	ctx := context.Background()
	// 设定租约，如果没有续期，就会自动删除
	lease := clientv3.NewLease(s.client)
	leaseRes, err := lease.Grant(ctx, 20)
	if err != nil {
		return err
	}
	_, err = kv.Put(context.Background(), key_prefix + id + "/" + name, addr, clientv3.WithLease(leaseRes.ID))
	if err != nil {
		return err
	}

	keepaliveRes, err :=lease.KeepAlive(context.Background(), leaseRes.ID)
	if err != nil {
		return err
	}

	go isKeepAlive(keepaliveRes)

	return err
}

func isKeepAlive(keepaliveRes <-chan *clientv3.LeaseKeepAliveResponse) {

	for {
		select {
		case ret := <-keepaliveRes:
			if ret != nil {
				fmt.Println("续租成功", time.Now())
			}
		}
	}

}


func (s *EtcdClientService) UnRegService(id string) error {
	kv := clientv3.NewKV(s.client)
	key_prefix := "/service/" + id
	_, err := kv.Delete(context.Background(), key_prefix, clientv3.WithPrefix())
	return err
}

