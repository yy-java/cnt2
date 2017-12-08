package user

import (
	"testing"

	. "github.com/yy-java/cnt2/db"

	"github.com/astaxie/beego/orm"
	. "github.com/smartystreets/goconvey/convey"
)

func init() {
	// set default database
	orm.RegisterDataBase("default", "mysql", "cnt2_db_user:q0NUVMca1@tcp(58.215.143.133:6307)/cnt2_db?charset=utf8", 30)
	orm.Debug = true
}

func TestSaveUserAuth(t *testing.T) {
	var uid int64 = 970916964
	app := "test app"

	Convey("Subject: Test save user auth", t, func() {
		auth := UserAuth{Uid: uid, Uname: "Chris", App: app, Permission: 9}
		SaveUserAuth(&auth)

		Convey("Test update operation", func() {
			auth.Uname = "Test"
			SaveUserAuth(&auth)

			authInDb, _ := FindUserAuthById(auth.Id)
			So(authInDb.Uname, ShouldEqual, "Test")

			RemoveUserAuthById(authInDb.Id)
		})

	})
}

func TestCheckManagePermission(t *testing.T) {
	var uid int64 = 970916964
	app := "test app"

	Convey("Subject: Test check user permission is manage", t, func() {
		SaveUserAuth(&UserAuth{Uid: uid, Uname: "Chris", App: app, Permission: 9})

		ismanager := CheckManagePermission(uid, app)
		So(ismanager, ShouldBeTrue)

		RemoveUserAuthByUidAndApp(uid, app)
	})
}
