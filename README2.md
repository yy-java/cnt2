### Initial Project
```
git clone http://github.com/yy-java/cnt2.git

```

### inital dependency
```
go get github.com/astaxie/beego //beego
go get github.com/beego/bee //beego tools
go get github.com/coreos/etcd/clientv3  // etcd client
go get github.com/go-sql-driver/mysql  // mysql driver
go install github.com/boltdb/bolt
go get github.com/coreos/etcd

go get github.com/smartystreets/goconvey // test
```

### Generate gRpc Code
```
go get -u github.com/golang/protobuf/protoc-gen-go
protoc -I pb/ pb/cnt2.proto --go_out=plugins=grpc:pb
```

### Generate thrift code
```
go get git.apache.org/thrift.git/lib/go/thrift
thrift -r --gen go secuserinfo.thrift
```

## project component
### grpcserver
API module to exposing gRPC service to SDKs

本地启动:
```
cd ${workspace}/github.com/yy-java/cnt2/grpcserver
go install
cd $GOPATH/bin/
bee run yy.com/cnt2/grpcserver
```
http://localhost:8080/v1/object

测试环境部署启动:
```
- 本地编译
cd ${workspace}/github.com/yy-java/cnt2/grpcserver
bee pack -be="GOOS=linux" com/cnt2/grpcserver

- 构建发布包
登陆http://yydeploy1.sysop.duowan.com/package/index.jspx?name=cnt2_grpcserver_test
创建cnt2_grpcserver_test包的新版本，填写新版本号
上传编译包${workspace}/github.com/yy-java\cnt2\grpcserver\grpcserver.tar.gz到发布包的$INSTALL_PATH/bin目录
保存打包后安装到
- 启动命令
sudo /data/services/cnt2_grpcserver_test-$VERSION/admin/start.sh
- 日志
cat /data/services/cnt2_grpcserver_test-$VERSION/admin/start.log
```

生产环境部署启动:
```
- 本地编译
cd ${workspace}/github.com/yy-java/cnt2/grpcserver
bee pack -be="GOOS=linux" com/cnt2/grpcserver

- 构建发布包
登陆http://yydeploy1.sysop.duowan.com/package/index.jspx?name=cnt2_grpcserver
创建cnt2_grpcserver包的新版本，填写新版本号
上传编译包${workspace}/github.com/yy-java\cnt2\grpcserver\grpcserver.tar.gz到发布包的$INSTALL_PATH/bin目录
保存打包后安装到221.228.91.155/221.228.83.119/113.108.65.32
- 启动命令
sudo /data/services/cnt2_grpcserver-$VERSION/admin/start.sh
- 日志
cat /data/services/cnt2_grpcserver-$VERSION/admin/start.log
```

### httpserver
API module to exposing http service to cnt console

本地启动:  
```
cd ${workspace}/github.com/yy-java/cnt2/httpserver
go install
cd $GOPATH/bin/
bee run yy.com/cnt2/httpserver
```
http://localhost:8081/app/test


测试环境部署启动:
```
- 本地编译
cd ${workspace}/github.com/yy-java/cnt2/httpserver
bee pack -be="GOOS=linux" com/cnt2/httpserver

- 构建发布包
登陆http://yydeploy1.sysop.duowan.com/package/index.jspx?name=cnt2_httpserver_test
创建cnt2_httpserver_test包的新版本，填写新版本号
上传编译包${workspace}/github.com/yy-java\cnt2\httpserver\httpserver.tar.gz到发布包的$INSTALL_PATH/bin目录
保存打包后安装到
- 启动命令
sudo /data/services/cnt2_httpserver_test-$VERSION/admin/start.sh
- 日志
cat /data/services/cnt2_httpserver_test-$VERSION/admin/start.log
- 探测地址
http://cnt-api-test.yy.com/app/test
```

生产环境部署启动:
```
- 本地编译
cd ${workspace}/github.com/yy-java/cnt2/httpserver
bee pack -be="GOOS=linux" com/cnt2/httpserver

- 构建发布包
登陆http://yydeploy1.sysop.duowan.com/package/index.jspx?name=cnt2_httpserver
创建cnt2_httpserver包的新版本，填写新版本号
上传编译包${workspace}/github.com/yy-java\cnt2\httpserver\httpserver.tar.gz到发布包的$INSTALL_PATH/bin目录
保存打包后安装到
- 启动命令
sudo /data/services/cnt2_httpserver-$VERSION/admin/start.sh
- 日志
cat /data/services/cnt2_httpserver-$VERSION/admin/start.log
- 探测地址
http://cnt-api.yy.com/app/test
```

- 安装使用命令

go get 
