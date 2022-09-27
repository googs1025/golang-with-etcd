# golang-with-etcd
## 基于etcd实现服务发现与服务注册。
## 其中支持多实例注册，client端采用随机负载均衡方法获取实例。
## 并且在server 被kill调时，实现优雅退出
