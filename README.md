# golang-with-etcd
```bigquery
.
├── README.md
├── bash.sh  // server 执行bash文件
├── client.go // 客户端
├── confdfiles // confd 配置文件
│   
├── lib // 服务发现与服务注册逻辑
│   ├── LoadBalance.go
│   ├── client.go
│   ├── endpoint.go
│   └── service.go
├── main.go
├── my.ini  // 配置文件
├── myserver.go  // 配置的server
└── service 
    ├── productRequest.go
    └── productTransports.go

```
## 基于etcd实现服务发现与服务注册。
[](https://github.com/googs1025/golang-with-etcd/blob/main/image/82192dd67a584bc9828e24e9f27f4a96.png?ram=true)
#### 1.本项目支持多实例注册，分为server端与client端。
#### 2.client端采用随机负载均衡方法获取实例。
#### 3.server端实现多实例注册，同时写入etcd，且定期续期，另外，也实现进程被kill掉时，优雅从etcd注销。
####项目启动步骤
#### 1. 启动server端
```bigquery
bash bash.sh 

(base) xxxxxxMacBook-Pro:etcd-practice zhenyu.jiang$ bash bash.sh 
续租成功 2022-09-29 21:04:44.013871 +0800 CST m=+0.349606376
续租成功 2022-09-29 21:04:44.075001 +0800 CST m=+0.230567001
续租成功 2022-09-29 21:04:44.425022 +0800 CST m=+0.346569417
续租成功 2022-09-29 21:04:50.988815 +0800 CST m=+7.324839960
续租成功 2022-09-29 21:04:51.08358 +0800 CST m=+7.239436751
```
#### 2. 启动client端
```bigquery
go run client.go
```

### 3. 进入etcd环境中查看
```bigquery
(base) xxxxxxMacBook-Pro:~ zhenyu.jiang$ etcdctl get /service --prefix
/service/3970ccb2-2ea8-4f93-a949-2eaa9678c4c1/productservice
127.0.0.1:1125
/service/82afbbcf-e6a6-43be-a8ca-a78a8a6c53b5/productservice
127.0.0.1:1124
/service/c0b00dc2-08ab-4c60-8b63-604312851ae9/productservice
127.0.0.1:1123
```



## 基于etcd实现配置监听与热更新
[](https://github.com/googs1025/golang-with-etcd/blob/main/image/%E6%B5%81%E7%A8%8B%E5%9B%BE%20(1).jpg?raw=true)
#### 项目步骤
#### 1. 查看my.ini文件
```bigquery
[db]
db_user = jiang
db_pass = 2334522233322
```

#### 2. 开启server端，监听8002端口
```bigquery
go run myserver.go -p 8002
```

#### 3. 修改my.ini文件
```bigquery
[db]
db_user = jiangjiang
db_pass = 111
```
### 4. 查看路由
```bigquery
在浏览器或postman上，执行http请求。
如：127.0.0.1:8002，即可显示。
127.0.0.1:8002/reload 会重新加载配置结果。
```

### 项目实现思路
```bigquery
1. 搭建server，写请求路由
2. 启goroutine，执行平滑重启。
3. 启goroutine，执行优雅退出机制
4. 启goroutine，监听配置文件变化，采用md5值的方式执行，只要值改变，就restart server。
```

## 3.基于etcd与confd实现配置更新
未更新完。。。


