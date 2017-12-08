package db

import (
	"log"

	"github.com/astaxie/beego/orm"
)

/** for table 'app' */
func (app *App) Create() error {
	o := orm.NewOrm()
	_, err := o.Insert(app)
	if err != nil {
		log.Printf("insert app {%v} failed, err: %v", app, err)
		return err
	}

	return nil
}

/** for table 'app' */
func (app *App) Update() (int64, error) {
	o := orm.NewOrm()

	num, err := o.Update(app)
	if err != nil {
		log.Printf("Update app {%v} failed, err: %v", app, err)
		return 0, err
	}

	return num, nil
}
func (app *App) Read() error {
	o := orm.NewOrm()
	err := o.Read(app)
	if err != nil {
		log.Printf("read app {%v} failed, err: %v", app, err)
		return err
	}
	return nil
}

func (app *App) Delete() (int64, error) {
	o := orm.NewOrm()
	num, err := o.Delete(app)
	if err != nil {
		log.Printf("delete app {%v} failed, err: %v", app, err)
		return 0, err
	}
	return num, nil
}

func (app *App) FindAll(uid int64) ([]*App, error) {
	o := orm.NewOrm()
	var apps []*App
	_, err := o.Raw("select distinct a.* from user_auth b,app a where b.uid=? and (a.app=b.app or b.permission=99) ", uid).QueryRows(&apps)
	if err != nil {
		log.Printf("FindAll uid {%v} failed, err: %v", uid, err)
		return nil, err
	}
	return apps, nil
}
