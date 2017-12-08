package db

import (
	"testing"

	"github.com/astaxie/beego/orm"
	. "github.com/smartystreets/goconvey/convey"
)

func init() {
	// set default database
	orm.RegisterDataBase(DB_Test_Name, DB_TEST_DRIVER, DB_TEST_DataSource, DB_TEST_MaxIdle)
	orm.Debug = true
}

func TestUserAuth(t *testing.T) {
	var uid int64 = 970916964
	app := "test app"

	Convey("Subject: Test table user_auth", t, func() {
		auth := UserAuth{Uid: uid, Uname: "Chris", App: app, Permission: 9}

		auth.Create()
		id := auth.Id

		Convey("Test read operation", func() {
			auth = UserAuth{Id: id}
			auth.Read()
			So(auth.Uid, ShouldEqual, uid)

			auth.Uname = "Test"
			num, _ := auth.Update()
			So(num, ShouldEqual, 1)

			auth = UserAuth{Id: id}
			auth.Read()
			So(auth.Uname, ShouldEqual, "Test")

			Convey("Test delete operation", func() {
				auth = UserAuth{Id: id}
				num, _ := auth.Delete()
				So(num, ShouldEqual, 1)
			})
		})
	})
}

func TestUserAuthCustomized(t *testing.T) {
	var uid int64 = 970916964
	app := "test app"

	Convey("Subject: Test table user_auth by customized operation", t, func() {
		auth := UserAuth{Uid: uid, Uname: "Chris", App: app, Permission: 9}

		auth.Create()

		Convey("Test read operation", func() {
			auth = UserAuth{Uid: uid, App: app}
			userAuths, _ := auth.ReadByInput()
			So(len(userAuths), ShouldEqual, 1)

			auth := userAuths[0]
			So(auth.Permission, ShouldEqual, 9)

			Convey("Test delete operation", func() {
				auth := UserAuth{Uid: uid, App: app}
				num, _ := auth.DeleteByInput()
				So(num, ShouldEqual, 1)
			})
		})
	})
}

func TestDeleteMultiUserAuth(t *testing.T) {
	Convey("test", t, func() {
		auth1 := UserAuth{Uid: 11111, App: "test"}
		auth1.Create()

		auth2 := UserAuth{Uid: 22222, App: "test"}
		auth2.Create()

		auth := UserAuth{App: "test"}
		auth.DeleteByInput()
	})
}
