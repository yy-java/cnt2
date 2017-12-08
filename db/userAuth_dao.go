package db

import (
	"log"

	"github.com/astaxie/beego/orm"
)

/** for table 'user_auth' */
func (auth *UserAuth) Create() error {
	o := orm.NewOrm()
	_, err := o.Insert(auth)
	if err != nil {
		log.Printf("insert user_auth {%v} failed, err: %v", auth, err)
		return err
	}

	return nil
}

func (auth *UserAuth) Update() (int64, error) {
	o := orm.NewOrm()
	num, err := o.Update(auth)
	if err != nil {
		log.Printf("update user_auth {%v} failed, err: %v", auth, err)
		return 0, err
	}
	return num, nil
}

func (auth *UserAuth) Read() error {
	o := orm.NewOrm()

	err := o.Read(auth)
	if err != nil {
		log.Printf("read user_auth {%v} failed, err: %v", auth, err)
		return err
	}
	return nil
}

func (auth *UserAuth) ReadByInput() ([]*UserAuth, error) {
	var userAuths []*UserAuth
	o := orm.NewOrm()

	num, err := buildUserAuthFilter(o, auth).All(&userAuths)
	if err != nil {
		log.Printf("read user_auth by input {%v} failed, err: %v", auth, err)
		return nil, err
	}
	if num == 0 {
		log.Printf("no record matched with user_auth {%v}, count: %d", auth, num)
		userAuths = make([]*UserAuth, 0)
	}

	return userAuths, nil
}

func (auth *UserAuth) Delete() (int64, error) {
	o := orm.NewOrm()
	num, err := o.Delete(auth)
	if err != nil {
		log.Printf("delete user_auth {%v} failed, err: %v", auth, err)
		return 0, err
	}
	return num, nil
}

func (auth *UserAuth) DeleteByInput() (int64, error) {
	o := orm.NewOrm()

	num, err := buildUserAuthFilter(o, auth).Delete()
	if err != nil {
		log.Printf("delete user_auth by input {%v} failed, err: %v", auth, err)
		return 0, err
	}
	return num, nil
}

func buildUserAuthFilter(o orm.Ormer, auth *UserAuth) orm.QuerySeter {
	qs := o.QueryTable(new(UserAuth))

	if auth.Id > 0 {
		qs = qs.Filter("Id", auth.Id)
	}
	if auth.Uid > 0 {
		qs = qs.Filter("Uid", auth.Uid)
	}
	if len(auth.Uname) > 0 {
		qs = qs.Filter("Uname", auth.Uname)
	}
	if len(auth.App) > 0 {
		qs = qs.Filter("App", auth.App)
	}
	if auth.Permission != int8(PermissionType_None) {
		qs = qs.Filter("Permission", auth.Permission)
	}

	return qs
}
