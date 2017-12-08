package db

import (
	"log"

	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
)

// Create create a node record and return node.Id
func (n *Node) Create() error {
	o := orm.NewOrm()
	_, err := o.Insert(n)
	if err != nil {
		log.Printf("insert Node {%v} failed, err: %v", n, err)
		return err
	}
	return nil
}

func (n *Node) Read() error {
	o := orm.NewOrm()
	err := o.Read(n)
	if err != nil {
		return err
	}
	return nil
}

func (n *Node) Delete() (int64, error) {
	o := orm.NewOrm()
	num, err := o.Delete(n)
	return num, err
}
func (n *Node) FindAllNode(app string, profile string) ([]*Node, error) {
	o := orm.NewOrm()
	var nodes []*Node
	_, err := o.QueryTable(new(Node)).Filter("app", app).Filter("profile", profile).All(&nodes)
	if err != nil {
		log.Printf("find node by {app: %s, profile: %s} failed, err: %v", app, profile, err)
		return nil, err
	}
	return nodes, nil
}
