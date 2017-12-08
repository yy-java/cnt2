package db

import (
	"github.com/astaxie/beego/orm"
	"log"
)

func (config *Config) Create() error {
	o := orm.NewOrm()
	_, err := o.Insert(config)
	if err != nil {
		log.Printf("insert config {%v} failed, err: %v", config, err)
		return err
	}

	return nil
}
func (config *Config) CreateOrUpdate() error {
	o := orm.NewOrm()
	_, err := o.InsertOrUpdate(config)
	if err != nil {
		log.Printf("insert config {%v} failed, err: %v", config, err)
		return err
	}

	return nil
}
func (config *Config) Update() (int64, error) {
	o := orm.NewOrm()
	num, err := o.Update(config)
	if err != nil {
		log.Printf("update config {%v} failed, err: %v", config, err)
		return 0, err
	}
	return num, nil
}

func (config *Config) Read() error {
	o := orm.NewOrm()

	err := o.Read(config)
	if err != nil {
		log.Printf("read config {%v} failed, err: %v", config, err)
		return err
	}
	return nil
}

func (config *Config) ReadByInput() ([]*Config, error) {
	var configs []*Config
	o := orm.NewOrm()

	num, err := buildConfigFilter(o, config).All(&configs)
	if err != nil {
		log.Printf("read config by input {%v} failed, err: %v", config, err)
		return nil, err
	}
	if num == 0 {
		log.Printf("no record matched with config {%v}, count: %d", config, num)
		return nil, nil
	}

	return configs, nil
}
func (config *Config) ReadAppProfiles(app string) ([]string, error) {
	o := orm.NewOrm()
	var profiles []string
	_, err := o.Raw("select distinct profile from `config` where app=?", app).QueryRows(&profiles)
	if err != nil {
		log.Printf("ReadAppProfiles app {%v} failed, err: %v", app, err)
		return nil, err
	}
	return profiles, nil
}

func (config *Config) Delete() (int64, error) {
	o := orm.NewOrm()
	num, err := o.Delete(config)
	if err != nil {
		log.Printf("delete config {%v} failed, err: %v", config, err)
		return 0, err
	}
	return num, nil
}

func (config *Config) DeleteByInput() (int64, error) {
	o := orm.NewOrm()

	num, err := buildConfigFilter(o, config).Delete()
	if err != nil {
		log.Printf("delete config by input {%v} failed, err: %v", config, err)
		return 0, err
	}
	return num, nil
}

func (config *Config) Count() (int64, error) {
	o := orm.NewOrm()

	count, err := buildConfigFilter(o, config).Count()
	if err != nil {
		log.Printf("count config by input {%v} failed, err: %v", config, err)
		return 0, err
	}
	return count, nil
}

func buildConfigFilter(o orm.Ormer, config *Config) orm.QuerySeter {
	qs := o.QueryTable(new(Config))

	if config.Id > 0 {
		qs = qs.Filter("id", config.Id)
	}
	if len(config.App) > 0 {
		qs = qs.Filter("app", config.App)
	}
	if len(config.Profile) > 0 {
		qs = qs.Filter("profile", config.Profile)
	}
	if len(config.Key) > 0 {
		qs = qs.Filter("key", config.Key)
	}
	if len(config.Value) > 0 {
		qs = qs.Filter("value", config.Value)
	}
	if config.Version > 0 {
		qs = qs.Filter("version", config.Version)
	}
	if len(config.PublishedValue) > 0 {
		qs = qs.Filter("published_value", config.PublishedValue)
	}
	if config.PublishedVersion > 0 {
		qs = qs.Filter("published_version", config.PublishedVersion)
	}

	return qs
}
func (config *Config) UpdateConfigPublishedValue(publishVersion int64, publishValue string) (int64, error) {
	o := orm.NewOrm()
	num, err := buildConfigFilter(o, config).Update(orm.Params{
		"PublishedValue":   publishValue,
		"PublishedVersion": publishVersion,
	})
	if err != nil {
		log.Printf("update config by input {%v} failed, err: %v", config, err)
		return 0, err
	}
	return num, nil
}

func (config *Config) CopyTo(destProfile, creator string, isManage bool) error {
	o := orm.NewOrm()
	approveType := 2
	approver := creator
	if !isManage {
		approveType = 1
		approver = ""
	}
	_, err := o.Raw("insert into config(app,profile,`key`,`value`,version,validator,description,modifier,modify_time,approve_type,approver) "+
		"   select app,?,`key`,`value`,1,validator,description,?,now(),?,? from config where app=? and `profile`=?", destProfile, creator, approveType, approver, config.App, config.Profile).Exec()
	if err != nil {
		log.Printf("CopyTo: {%v}  to destProfile {%v} failed, err: %v", config, destProfile, err)
		return err
	}
	return nil
}
