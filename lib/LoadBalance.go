package lib

import (
	"math/rand"
	"time"
)

type LoadBalance struct {
	Servers []*ServiceInfo
}

func NewloadBalance(servers []*ServiceInfo) *LoadBalance  {
	return &LoadBalance{Servers:servers}
}
func(l *LoadBalance) getByRand(sname string ) *ServiceInfo{
	tmp:=make([]*ServiceInfo,0)
	for _,service:=range l.Servers{
		if service.ServiceName==sname{
			tmp=append(tmp,service)
		}
	}
	if len(tmp)==0 {
		return nil
	}
	rand.Seed(time.Now().UnixNano())
	i:=rand.Intn(len(tmp))
	return l.Servers[i]

}
