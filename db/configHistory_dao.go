package db

import (
	"log"

	"github.com/astaxie/beego/orm"
)

/** for table 'config_history' */
func (ch *ConfigHistory) Create() error {
	o := orm.NewOrm()
	_, err := o.Insert(ch)
	if err != nil {
		log.Printf("insert config_history {%v} failed, err: %v", &ch, err)
		return err
	}

	return nil
}

// 通过主键更新所有值
func (ch *ConfigHistory) Update() (int64, error) {
	o := orm.NewOrm()
	num, err := o.Update(ch)
	if err != nil {
		log.Printf("update config_history {%v} failed, err: %v", &ch, err)
		return 0, err
	}
	return num, nil
}

// 通过主键查询
func (ch *ConfigHistory) Read() (*ConfigHistory, error) {
	o := orm.NewOrm()

	err := o.Read(ch)
	if err != nil {
		log.Printf("read config_history {id: %d} failed, err: %v", ch.Id, err)
		return nil, err
	}
	return ch, nil
}

// 通过各个有指定的参数查询
func (ch *ConfigHistory) ReadByInput() ([]*ConfigHistory, error) {
	var configHistorys []*ConfigHistory
	o := orm.NewOrm()

	num, err := buildConfigHistoryFilter(o, ch).All(&configHistorys)
	if err != nil {
		log.Printf("read config_history by input {%v} failed, err: %v", &ch, err)
		return nil, err
	}
	if num == 0 {
		log.Printf("no record matched with config_history {%v}, count: %d", &ch, num)
		return nil, nil
	}

	return configHistorys, nil
}

func (ch *ConfigHistory) ReadMaxVersionHistory() (*ConfigHistory, error) {
	var configHistorys []*ConfigHistory

	o := orm.NewOrm()
	qs := o.QueryTable(new(ConfigHistory))

	num, err := qs.Filter("app", ch.App).Filter("profile", ch.Profile).Filter("key", ch.Key).OrderBy("-cur_version").Limit(1).All(&configHistorys)
	if err != nil {
		log.Printf("read max version of config_history by {app: %s, profiile: %s, key: %s} failed, err: %v", ch.App, ch.Profile, ch.Key, err)
		return nil, err
	}
	if num != 1 {
		log.Printf("no record matched with config_history {app: %s, profiile: %s, key: %s}, count: %d", ch.App, ch.Profile, ch.Key, num)
		return nil, nil
	}

	return configHistorys[0], nil
}

// 通过主键删除
func (ch *ConfigHistory) Delete() (int64, error) {
	o := orm.NewOrm()
	num, err := o.Delete(ch)
	if err != nil {
		log.Printf("delete config_history {id: %d} failed, err: %v", ch.Id, err)
		return 0, err
	}
	return num, nil
}

// 通过有指定的参数作为where条件删除（note：会先执行select，再执行delete）
func (ch *ConfigHistory) DeleteByInput() (int64, error) {
	o := orm.NewOrm()

	num, err := buildConfigHistoryFilter(o, ch).Delete()
	if err != nil {
		log.Printf("delete config_history by input {%v} failed, err: %v", &ch, err)
		return 0, err
	}
	return num, nil
}

func buildConfigHistoryFilter(o orm.Ormer, configHistory *ConfigHistory) orm.QuerySeter {
	qs := o.QueryTable(new(ConfigHistory))

	if configHistory.Id > 0 {
		qs = qs.Filter("id", configHistory.Id)
	}
	if len(configHistory.App) > 0 {
		qs = qs.Filter("app", configHistory.App)
	}
	if len(configHistory.Profile) > 0 {
		qs = qs.Filter("profile", configHistory.Profile)
	}
	if len(configHistory.Key) > 0 {
		qs = qs.Filter("key", configHistory.Key)
	}
	if configHistory.CurVersion > 0 {
		qs = qs.Filter("cur_version", configHistory.CurVersion)
	}
	qs = qs.OrderBy("-id").Limit(20)

	return qs
}
