package main

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"log"
	_ "github.com/yy-java/cnt2/httpserver/routers"
)

func main() {
	

	log.Println("cnt2 runmode", beego.BConfig.RunMode)
	log.Println("etcd::endpoints", beego.AppConfig.Strings("etcd.endpoints"))
	log.Println("mysql::dataSource", beego.AppConfig.Strings("mysql.dataSource"))

	orm.RegisterDataBase("default", "mysql", beego.AppConfig.String("mysql.dataSource"), 30, 100)

	beego.Run()
}
