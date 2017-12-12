package main

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"log"
	_ "github.com/yy-java/cnt2/httpserver/routers"
)

func main() {

	if beego.BConfig.RunMode == "dev" {
		beego.BConfig.WebConfig.DirectoryIndex = true
		beego.BConfig.WebConfig.StaticDir["/swagger"] = "swagger"

		beego.AppConfig.Set("etcd::endpoints", "127.0.0.1:2379")
		beego.AppConfig.Set("mysql::dataSource", "username:pwd@tcp(127.0.0.1:3306)/cnt2_db?charset=utf8&loc=Asia%2FShanghai")

		orm.Debug = true
	}

	log.Println("cnt2 runmode", beego.BConfig.RunMode)
	log.Println("etcd::endpoints", beego.AppConfig.Strings("etcd::endpoints"))
	log.Println("mysql::dataSource", beego.AppConfig.Strings("mysql::dataSource"))

	orm.RegisterDataBase("default", "mysql", beego.AppConfig.String("mysql::dataSource"), 30, 100)

	beego.Run()
}
