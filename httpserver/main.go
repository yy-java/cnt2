package main

import (
	//	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"log"
	//	"os"
	//	"time"
	_ "github.com/yy-java/cnt2/httpserver/routers"
	"github.com/yy-java/cnt2/service/account"
)

func main() {
	//	log.SetFlags(log.LstdFlags | log.Lshortfile)
	//	logFile, err := os.OpenFile("/data/weblog/business/cnt2_httpserver/web.log", os.O_RDWR|os.O_CREATE, 0666)
	//	if err != nil {
	//		fmt.Printf("open file error=%s\r\n", err.Error())
	//		os.Exit(-1)
	//	}
	//	log.SetOutput(logFile)
	//	log.Println("starting httpserver at", time.Now().Format("2006-01-02 15:04:05"))

	if beego.BConfig.RunMode == "dev" {
		beego.BConfig.WebConfig.DirectoryIndex = true
		beego.BConfig.WebConfig.StaticDir["/swagger"] = "swagger"

		beego.AppConfig.Set("etcd::endpoints", "61.147.187.152:2379;61.147.187.142:2379;61.147.187.150:2379")
		beego.AppConfig.Set("mysql::dataSource", "cnt2_db_user:q0NUVMca1@tcp(58.215.143.133:6307)/cnt2_db?charset=utf8&loc=Asia%2FShanghai")

		orm.Debug = true
	}

	log.Println("cnt2 runmode", beego.BConfig.RunMode)
	log.Println("etcd::endpoints", beego.AppConfig.Strings("etcd::endpoints"))
	log.Println("mysql::dataSource", beego.AppConfig.Strings("mysql::dataSource"))

	orm.RegisterDataBase("default", "mysql", beego.AppConfig.String("mysql::dataSource"), 30, 100)
	account.InitUdbCookieService(beego.AppConfig.String("UDB::appId"), beego.AppConfig.String("UDB::appKey"))

	beego.Run()
}
