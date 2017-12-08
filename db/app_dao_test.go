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

func TestApp(t *testing.T) {

	Convey("Subject: Test table:app operation", t, func() {
		app := App{App: "test app", AppType: 0, Name: "test app", Charger: "george", CreateTime: time.Now()}

		app.Create()

		Convey("Test read operation", func() {
			app = App{App: "test app"}
			app.Read()
			So(app.Name, ShouldEqual, "test app")

			Convey("Test delete operation", func() {
				num, _ := app.Delete()
				So(num, ShouldEqual, 1)
			})
		})

	})
}
