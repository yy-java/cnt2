package db

import (
	"testing"
	"time"

	"github.com/astaxie/beego/orm"
	. "github.com/smartystreets/goconvey/convey"
)

func init() {
	// set default database
	orm.RegisterDataBase(DB_Test_Name, DB_TEST_DRIVER, DB_TEST_DataSource, DB_TEST_MaxIdle)
	orm.Debug = true
}

func Test_NodeCreate(t *testing.T) {
	Convey("Subject: Test table:node operation", t, func() {
		node := Node{App: "test", Profile: "WX", Sip: "127.0.0.1", Pid: 123, RegisterTime: time.Now()}

		node.Create()
		id := node.Id

		Convey("Test read operation", func() {
			node = Node{Id: id}
			node.Read()
			So(node.Id, ShouldEqual, id)

			Convey("Test delete operation", func() {
				node = Node{Id: id}
				num, _ := node.Delete()
				So(num, ShouldEqual, 1)
			})
		})

	})
}
