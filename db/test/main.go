package main

import (
	"github.com/yy-java/cnt2/db"
	
	"github.com/astaxie/beego/orm"
)

const (
	DB_Test_Name       = "default"
	DB_TEST_DRIVER     = "mysql"
	DB_TEST_DataSource = "cnt2_db_user:q0NUVMca1@tcp(58.215.143.133:6307)/cnt2_db?charset=utf8"
	DB_TEST_MaxIdle    = 30
)

func init() {
	// set default database
	orm.RegisterDataBase(DB_Test_Name, DB_TEST_DRIVER, DB_TEST_DataSource, DB_TEST_MaxIdle)
	orm.Debug = true
}

func main() {
	nodePublish := db.NodePublish{Id:0, App:"app", PublishResult:-1}

	nodePublish.UpdatePublishResult()
}