package db

import (
	"log"
	"github.com/astaxie/beego/orm"
)

/** for table 'user' */
func (user *User) Create() error {
	o := orm.NewOrm()
	_, err := o.Insert(user)
	if err != nil {
		log.Printf("insert user {%v} failed, err: %v", user, err)
		return err
	}
	return nil
}
 
func (user *User) Read() error {
	o := orm.NewOrm()
	err := o.Read(user)
	if err != nil {
		log.Printf("read user {%v} failed, err: %v", user, err)
		return err
	}
	return nil
}

func (user *User) Delete() (int64, error) {
	o := orm.NewOrm()
	num, err := o.Delete(user)
	if err != nil {
		log.Printf("delete user {%v} failed, err: %v", user, err)
		return 0, err
	}
	return num, nil
}

func (user *User) FindAll() ([]*User, error) {
	o := orm.NewOrm()
	var apps []*User
	_, err := o.Raw("select * from user order by create_time desc").QueryRows(&apps)
	if err != nil {
		log.Printf("FindAll failed, err: %v", err)
		return nil, err
	}
	return apps, nil
}
func (user *User) Login() (*User, error) {
	o := orm.NewOrm()
	var u  *User
	_, err := o.Raw("select * from user where username =? and pwd=?  limit 1",user.Username,user.Pwd).QueryRows(&u)
	if err != nil {
		log.Printf("Login username {%v} failed, err: %v", user.Username,err)
		return nil, err
	}
	return u, nil
}
