package db

import (
	"log"
	"testing"

	"github.com/astaxie/beego/orm"
	. "github.com/smartystreets/goconvey/convey"
)

func init() {
	// set default database
	orm.RegisterDataBase(DB_Test_Name, DB_TEST_DRIVER, DB_TEST_DataSource, DB_TEST_MaxIdle)
	orm.Debug = true
}

func TestNodePublish(t *testing.T) {
	var nodeId int64 = 2
	app := "app1"

	Convey("Subject: Test table user_auth", t, func() {
		nodePublish := NodePublish{NodeId: nodeId, App: app, Profile: "dev", Key: "key2"}

		nodePublish.Create()
		id := nodePublish.Id

		log.Println(id)

		Convey("Test read operation", func() {
			nodePublish, _ := nodePublish.Read()
			So(nodePublish.NodeId, ShouldEqual, nodeId)

			Convey("Test delete operation", func() {
				num, _ := nodePublish.Delete()
				So(num, ShouldEqual, 1)
			})
		})
	})
}
