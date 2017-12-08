package db

import (
	"github.com/astaxie/beego/orm"
	"log"
)

/** for table 'profile' */
func (profile *Profile) Create() error {
	o := orm.NewOrm()
	_, err := o.Insert(profile)
	if err != nil {
		log.Printf("insert profile {%v} failed, err: %v", profile, err)
		return err
	}

	return nil
}

/** for table 'profile' */
func (profile *Profile) Update() (int64, error) {
	o := orm.NewOrm()

	num, err := o.InsertOrUpdate(profile, "App", "Profile", "Name")
	if err != nil {
		log.Printf("Update profile {%v} failed, err: %v", profile, err)
		return 0, err
	}

	return num, nil
}
func (profile *Profile) Read() error {
	o := orm.NewOrm()
	err := o.Read(profile, "App", "Profile")
	if err != nil {
		log.Printf("read profile {%v} failed, err: %v", profile, err)
		return err
	}
	return nil
}

func (profile *Profile) Delete() (int64, error) {
	o := orm.NewOrm()
	num, err := o.Delete(profile, "App", "Profile")
	if err != nil {
		log.Printf("delete profile {%v} failed, err: %v", profile, err)
		return 0, err
	}
	return num, nil
}
func (profile *Profile) DeleteByApp() (int64, error) {
	o := orm.NewOrm()
	_, err := o.Raw("delete from profile where app=?", profile.App).Exec()
	if err != nil {
		log.Printf("DeleteByApp by input {%v} failed, err: %v", profile, err)
		return 0, err
	}
	return 1, nil
}
func (profile *Profile) DeleteByAppAndProfile() (int64, error) {
	o := orm.NewOrm()
	_, err := o.Raw("delete from profile where app=? and profile=?", profile.App, profile.Profile).Exec()
	if err != nil {
		log.Printf("DeleteByAppAndProfile by input {%v} failed, err: %v", profile, err)
		return 0, err
	}
	return 1, nil
}

func (profile *Profile) FindByApp(app string) ([]*Profile, error) {
	o := orm.NewOrm()
	var profiles []*Profile
	_, err := o.QueryTable(new(Profile)).Filter("app", app).All(&profiles)
	if err != nil {
		log.Printf("find profile by {app: %s} failed, err: %v", app, err)
		return nil, err
	}
	return profiles, nil
}
