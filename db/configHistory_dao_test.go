package db

import (
	"testing"
	"log"

	"github.com/astaxie/beego/orm"
	. "github.com/smartystreets/goconvey/convey"
)

func init() {
	// set default database
	orm.RegisterDataBase(DB_Test_Name, DB_TEST_DRIVER, DB_TEST_DataSource, DB_TEST_MaxIdle)
	orm.Debug = true
}

func TestConfigHistory(t *testing.T) {
	app := "app1"

	Convey("Subject: Test table user_auth", t, func() {
		configHistory := ConfigHistory{App: app, Profile:"dev", Key:"key2"}

		configHistory.Create()
		id := configHistory.Id
		
		log.Println(id)

		Convey("Test read operation", func() {
			configHistory, _ := configHistory.Read()
			So(configHistory.App, ShouldEqual, app)

			Convey("Test delete operation", func() {
				num, _ := configHistory.Delete()
				So(num, ShouldEqual, 1)
			})
		})
	})
}