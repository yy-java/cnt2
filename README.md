[console]: https://github.com/yy-java/cnt2-console
[gosdk]: https://github.com/yy-java/cnt2-gosdk
[javasdk]: https://github.com/yy-java/cnt2-javasdk
[etcd]: https://github.com/coreos/etcd
[etcd-cluster-install]: https://github.com/coreos/etcd/blob/master/Documentation/op-guide/clustering.md
[go-install]: https://golang.org/dl/

# cnt2

   cnt2是高可用的分布式配置中心, 是config center v2的简写， v2版本是在我们公司使用的v1版本上改进而来的。
采用go语言编写，客户端通过监听[etcd][etcd]集群感知配置的变化，再通过gRPC调用GrpcServer服务查询保存在
Mysql中最新的配置信息。

具有如下的功能:

* *实时性*: 一键发布，立马生效
* *多环境*: 自定义环境，一般分，开发、测试、生产环境1、生产环境2...生产环境N。多生产环境间接实现了分集群配置管理。
* *灰度发布*: 灰度发布，降低风险
* *配置回滚*: 能查看配置的历史版本，可以回滚到任意一个版本
* *审核机制*: 只能发布已审核的配置，出了问题拉上审核者一块来背锅^_^
* *权限管理*: 基于App的维度分管理员和开发者，还可以设置超级管理员来管理所有的App
* *后台管理*: 清新简洁的[后台管理][console]，与HttpServer进程提供的接口交互。
* *客户端SDK*: 目前提供[java][javasdk] 和 [go][gosdk] 两种语言的SDK

## 项目结构图

![structure](statics/structure.png)

## 开始

### 安装go

 [go][go-install]版本必须是1.9+
 
### 安装etcd

 查看[etcd集群安装文档][etcd-cluster-install]

### 建立Mysql数据表
 通过mysql目录下cnt2_db.sql建库建表
 
 
### 安装必须的go依赖库
```
go get github.com/astaxie/beego //beego
go get github.com/beego/bee //beego tools
go get github.com/coreos/etcd/clientv3  // etcd client
go get github.com/go-sql-driver/mysql  // mysql driver
go get github.com/smartystreets/goconvey // test
```

### Download Code

```
cd ${download-path}

git clone https://github.com/yy-java/cnt2.git
```

### 启动GrpcServer

cnt2也是一个自给自足的服务治理中心，所以建议部署多个GrpcServer实例，java/go sdk会轮询访问这些实例。
```
cd ${download-path}/cnt2/grpcserver
go install
cd $GOPATH/bin/
bee run yy.com/cnt2/grpcserver
```

### 启动HttpServer

```
cd ${download-path}/cnt2/httpserver
go install
cd $GOPATH/bin/
bee run yy.com/cnt2/httpserver
```

### 启动Consle


请查看[后台管理][console]

