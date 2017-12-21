package db

import (
	"time"
)

type PermissionType int8

const (
	PermissionType_None   PermissionType = 0
	PermissionType_Dev    PermissionType = 1
	PermissionType_Manage PermissionType = 9
	PermissionType_Admin  PermissionType = 99
)

type PublishResult int8

const (
	PublishResult_Create  PublishResult = 0
	PublishResult_Success PublishResult = 1
	PublishResult_Fail    PublishResult = 2
)

type PublishType int8

const (
	PublishType_Gray PublishType = 0
	PublishType_Full PublishType = 1
)

type ApproveType int8

const (
	ApproveType_NotApprove ApproveType = 1
	ApproveType_PASS       ApproveType = 2
)

type OperationType int8

const (
	OperationType_Create   OperationType = 1
	OperationType_Delete   OperationType = 2
	OperationType_Modify   OperationType = 3
	OperationType_Rollback OperationType = 4
)

type App struct {
	App        string    `orm:"pk;column(app)" json:"app"`
	AppType    int8      `orm:"column(app_type)" json:"appType"`
	Name       string    `orm:"size(100)" json:"name"`
	Charger    string    `orm:"size(50)" json:"charger"`
	CreateTime time.Time `orm:"auto_now;column(create_time);type(datetime)" json:"createTime"`
	ChargerUid int64     `orm:"column(charger_uid)" json:"chargerUid"`
}
type AppExt struct {
	App
	Permission int8 `json:"permission"`
}

type UserAuth struct {
	Id         int64  `orm:"pk;auto" json:"id"`
	Uid        int64  `orm:"column(uid)" json:"uid"`
	Uname      string `orm:"size(50)" json:"uname"`
	App        string `orm:"size(50)" json:"app"`
	Permission int8   `orm:"size(1)" json:"permission"`
}

type User struct {
	Uid         int64  `orm:"pk;auto" json:"uid"`
	Username      string `orm:"size(40)" json:"username"`
	Pwd        string `orm:"size(64)" json:"pwd"`
	CreateTime  time.Time `orm:"auto_now;column(create_time)" json:"createTime"`
}

type Config struct {
	Id               int64     `orm:"pk;column(id);auto" json:"id"`
	App              string    `orm:"column(app)" json:"app"`
	Profile          string    `orm:"column(profile)" json:"profile"`
	Key              string    `orm:"column(key)" json:"key"`
	Value            string    `orm:"column(value)" json:"value"`
	Version          int64     `orm:"column(version)" json:"version"`
	PublishedValue   string    `orm:"column(published_value)" json:"publishedValue"`
	PublishedVersion int64     `orm:"column(published_version)" json:"publishedVersion"`
	Validator        string    `orm:"column(validator)" json:"validator"`
	Modifier         string    `orm:"column(modifier)" json:"modifier"`
	ModifyTime       time.Time `orm:"auto_now;column(modify_time)" json:"modifyTime"`
	Description      string    `orm:"column(description)" json:"description"`
	Approver         string    `orm:"column(approver)" json:"approver"`
	ApproveType      int8      `orm:"column(approve_type)" json:"approveType"`
}

type ConfigHistory struct {
	Id          int64     `orm:"pk;column(id);auto" json:"id"`
	App         string    `orm:"column(app)" json:"app"`
	Profile     string    `orm:"column(profile)" json:"profile"`
	Key         string    `orm:"column(key)" json:"key"`
	PreValue    string    `orm:"column(pre_value)" json:"preValue"`
	CurValue    string    `orm:"column(cur_value)" json:"curValue"`
	PreVersion  int64     `orm:"column(pre_version)" json:"preVersion"`
	CurVersion  int64     `orm:"column(cur_version)" json:"curVersion"`
	Validator   string    `orm:"column(validator)" json:"validator"`
	Modifier    string    `orm:"column(modifier)" json:"modifier"`
	ModifyTime  time.Time `orm:"auto_now;column(modify_time)" json:"modifyTime"`
	Description string    `orm:"column(description)" json:"description"`
	OperateType int8      `orm:"column(operate_type)" json:"operateType"`
}

type Node struct {
	Id           int64     `orm:"pk;column(id);auto" json:"id"`
	App          string    `orm:"column(app)" json:"app"`
	Profile      string    `orm:"column(profile)" json:"profile"`
	Sip          string    `orm:"column(sip)" json:"sip"`
	Pid          int       `orm:"column(pid)" json:"pid"`
	RegisterTime time.Time `orm:"auto_now_add;column(register_time)" json:"registerTime"`
}

type NodePublish struct {
	Id            int64     `orm:"pk;column(id);auto" json:"id"`
	NodeId        int64     `orm:"column(node_id)" json:"nodeId"`
	App           string    `orm:"column(app)" json:"app"`
	Profile       string    `orm:"column(profile)" json:"profile"`
	Key           string    `orm:"column(key)" json:"key"`
	Version       int64     `orm:"column(version)" json:"version"`
	PublishTime   time.Time `orm:"auto_now;column(publish_time)" json:"publishTime"`
	PublishResult int8      `orm:"column(publish_result)" json:"publishResult"`
	PublishType   int8      `orm:"column(publish_type)" json:"publishType"`
}
type NodePublishExt struct {
	NodePublish
	Sip string `orm:"column(sip)" json:"sip"`
}

type Profile struct {
	App        string    `orm:"pk;column(app)" json:"app"`
	Profile    string    `orm:"column(profile)" json:"profile"`
	Name       string    `orm:"column(name)" json:"name"`
	CreateTime time.Time `orm:"auto_now;column(create_time);type(datetime)" json:"createTime"`
}
