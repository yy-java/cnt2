package main

import (
	"github.com/yy-java/cnt2/etcd"

	"log"
	"time"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/yy-java/cnt2/grpc/server"
	"github.com/yy-java/cnt2/grpcserver/globals"
	"github.com/yy-java/cnt2/service/models"
	"github.com/yy-java/cnt2/service/utils"
)

func GrpcServerRegister(port int) {
	hi := utils.GetHostInfo("")
	log.Println("register to etcd hostInfo:%s", hi)
	if len(hi.Ips) == 0 {
		return // ????
	}
	var info models.GrpcServerRegisterInfo
	info.ServerIP = make(map[int]string)
	for _, ip := range hi.Ips {
		info.ServerIP[ip.Type.Val] = ip.Ip
	}
	info.GroupId = hi.GroupId
	info.Port = port
	j, err := info.ToJson()
	if err != nil {
		log.Fatalln("Parse GrpcServerRegisterInfo error", err.Error())
		return
	}
	key := etcd.GetGrpcServerKey(hi.Ips[0].Ip, port)
	log.Println("register to etcd key:%s, value:%s", key, j)
	globals.GetEtcdClient().Lease(key, j)
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.Println("starting grpcserver at", time.Now().Format("2006-01-02 15:04:05"))

	if beego.BConfig.RunMode == "dev" {
		beego.BConfig.WebConfig.DirectoryIndex = true
		beego.BConfig.WebConfig.StaticDir["/swagger"] = "swagger"

		beego.AppConfig.Set("etcd::endpoints", "172.27.141.11:2379;172.27.141.12:2379;172.27.141.13:2379")
		beego.AppConfig.Set("grpc:serverPort", "50051")
		beego.AppConfig.Set("mysql::dataSource", "cnt2_db_user:q0NUVMca1@tcp(58.215.143.133:6307)/cnt2_db?charset=utf8&loc=Asia%2FShanghai")

		orm.Debug = true
	}

	log.Println("cnt2 runmode", beego.BConfig.RunMode)
	log.Println("etcd::endpoints", beego.AppConfig.Strings("etcd::endpoints"))
	log.Println("grpc:serverPort", beego.AppConfig.Strings("grpc::serverPort"))
	log.Println("mysql::dataSource", beego.AppConfig.Strings("mysql::dataSource"))

	orm.RegisterDataBase("default", "mysql", beego.AppConfig.String("mysql::dataSource"), 30, 100)
	grpcServerPort, _ := beego.AppConfig.Int("grpc::serverPort")
	grpcServer := server.NewServer(grpcServerPort)
	GrpcServerRegister(grpcServerPort)
	server.Start(grpcServer)
	log.Println("start sucess.............................")
	out := make(chan int)
	<- out
}
