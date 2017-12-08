package db

import (
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
)

func init() {
	//dataSource := beego.AppConfig.String("mysql::dataSource")

	// register model
	orm.RegisterModel(new(App), new(UserAuth), new(Config), new(Node), new(ConfigHistory), new(NodePublish), new(Profile))

	// set default database
	//orm.RegisterDataBase("default", "mysql", dataSource, 30, 100)

	orm.Debug = true
}
