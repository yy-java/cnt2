package db

import (
	"log"

	"github.com/astaxie/beego/orm"
)

/** for table 'node_publish' */
func (np *NodePublish) Create() error {
	o := orm.NewOrm()
	_, err := o.Insert(np)
	if err != nil {
		log.Printf("insert node_publish {%v} failed, err: %v", &np, err)
		return err
	}

	return nil
}

// 通过主键更新所有值
func (np *NodePublish) Update() (int64, error) {
	o := orm.NewOrm()
	num, err := o.Update(np)
	if err != nil {
		log.Printf("update node_publish {%v} failed, err: %v", &np, err)
		return 0, err
	}
	return num, nil
}

// 通过主键更新publish_result
func (np *NodePublish) UpdatePublishResult() (int64, error) {
	o := orm.NewOrm()
	num, err := o.Update(np, "publish_result")
	if err != nil {
		log.Printf("update node_publish {%v} failed, err: %v", &np, err)
		return 0, err
	}
	return num, nil
}

// 通过主键查询
func (np *NodePublish) Read() (*NodePublish, error) {
	o := orm.NewOrm()

	err := o.Read(np)
	if err != nil {
		log.Printf("read node_publish {id: %d} failed, err: %v", np.Id, err)
		return nil, err
	}
	return np, nil
}

// 通过各个有指定的参数查询
func (np *NodePublish) ReadByInput() ([]*NodePublish, error) {
	var nodePublishs []*NodePublish
	o := orm.NewOrm()

	num, err := buildNodePublishFilter(o, np).All(&nodePublishs)
	if err != nil {
		log.Printf("read node_publish by input {%v} failed, err: %v", &np, err)
		return nil, err
	}
	if num == 0 {
		log.Printf("no record matched with node_publish {%v}, count: %d", &np, num)
		nodePublishs = make([]*NodePublish, 0)
	}

	return nodePublishs, nil
}

// 通过主键删除
func (np *NodePublish) Delete() (int64, error) {
	o := orm.NewOrm()
	num, err := o.Delete(np)
	if err != nil {
		log.Printf("delete node_publish {id: %d} failed, err: %v", np.Id, err)
		return 0, err
	}
	return num, nil
}

// 通过有指定的参数作为where条件删除（note：会先执行select，再执行delete）
func (np *NodePublish) DeleteByInput() (int64, error) {
	o := orm.NewOrm()

	num, err := buildNodePublishFilter(o, np).Delete()
	if err != nil {
		log.Printf("delete node_publish by input {%v} failed, err: %v", &np, err)
		return 0, err
	}
	return num, nil
}

func buildNodePublishFilter(o orm.Ormer, nodePublish *NodePublish) orm.QuerySeter {
	qs := o.QueryTable(new(NodePublish))

	if nodePublish.Id > 0 {
		qs = qs.Filter("id", nodePublish.Id)
	}
	if nodePublish.NodeId > 0 {
		qs = qs.Filter("nodeId", nodePublish.NodeId)
	}
	if len(nodePublish.App) > 0 {
		qs = qs.Filter("app", nodePublish.App)
	}
	if len(nodePublish.Profile) > 0 {
		qs = qs.Filter("profile", nodePublish.Profile)
	}
	if len(nodePublish.Key) > 0 {
		qs = qs.Filter("key", nodePublish.Key)
	}
	if nodePublish.Version > 0 {
		qs = qs.Filter("version", nodePublish.Version)
	}
	if nodePublish.PublishResult >= 0 {
		qs = qs.Filter("publish_result", nodePublish.PublishResult)
	}
	if nodePublish.PublishType > 0 {
		qs = qs.Filter("publish_type", nodePublish.PublishType)
	}

	return qs
}

func (np *NodePublish) QueryPublishHistory(app string, profile string, key string, version int64) ([]*NodePublishExt, error) {
	o := orm.NewOrm()
	var result []*NodePublishExt
	_, err := o.Raw("select a.sip,b.version, b.publish_result,b.publish_time,b.publish_type from node a , node_publish b where a.id = b.node_id and b.app=? and b.profile=? and b.key=? and b.version =?", app, profile, key, version).QueryRows(&result)
	if err != nil {
		log.Printf("QueryPublishHistory app {%v} failed, err: %v", app, err)
		return nil, err
	}
	return result, nil
}

func (np *NodePublish) FindPublishedNode(app, profile, key string, version int64) ([]string, error) {
	o := orm.NewOrm()
	var result []string
	_, err := o.Raw("select node_id from node_publish where app=? and profile in (?,'common') and `key`=? and version =?", app, profile, key, version).QueryRows(&result)
	if err != nil {
		log.Printf("FindPublishedNode app {%v} failed, err: %v", app, err)
		return nil, err
	}
	return result, nil
}
