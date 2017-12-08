package db

import (
	"testing"

	"github.com/astaxie/beego/orm"
	. "github.com/smartystreets/goconvey/convey"
	"time"
)

func init() {
	// set default database
	orm.RegisterDataBase(DB_Test_Name, DB_TEST_DRIVER, DB_TEST_DataSource, DB_TEST_MaxIdle)
	orm.Debug = true
}

func TestConfigOperation(t *testing.T) {
	Convey("Subject: Test Config operation", t, func() {
		config := Config{App: "Test", Profile: "Dev", Key: "foo", Value: "bar", Version: 1, Modifier: "Chris", ModifyTime: time.Now(), Description: "for test"}

		config.Create()

		Convey("Test RUD operation", func() {
			id := config.Id
			newConfig := Config{Id: id}
			newConfig.Read()
			So(newConfig.Key, ShouldEqual, "foo")

			newConfig.Key = "testKey"
			newConfig.Value = "testVal"
			newConfig.Update()

			newConfig = Config{Id: id}
			newConfig.Read()
			So(newConfig.Key, ShouldEqual, "testKey")

			num, _ := newConfig.Delete()
			So(num, ShouldEqual, 1)
		})

		Convey("Test cusomized operation", func() {
			newConfig := Config{Key: "foo", Value: "bar"}
			configs, _ := newConfig.ReadByInput()
			So(len(configs), ShouldEqual, 1)

			count, _ := newConfig.Count()
			So(count, ShouldEqual, 1)

			num, _ := newConfig.DeleteByInput()
			So(num, ShouldEqual, 1)

			count, _ = newConfig.Count()
			So(count, ShouldEqual, 0)
		})
	})
}
