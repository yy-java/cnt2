package config

import (
	"log"

	"strconv"
	"time"
	"github.com/yy-java/cnt2/db"
	. "github.com/yy-java/cnt2/service/errors"
	"github.com/yy-java/cnt2/service/user"

	"github.com/yy-java/cnt2/etcd"
	"github.com/yy-java/cnt2/httpserver/globals"
)

const (
	Init_Version int64 = 1
)

func FindAllConfig(app string, profile string) ([]*db.Config, error) {
	if len(app) == 0 || len(profile) == 0 {
		return nil, ErrInvalidParam
	}

	cf := db.Config{App: app, Profile: profile}
	cfs, err := cf.ReadByInput()

	if err != nil {
		log.Printf("find configs by {app: %s, profile: %s} failed, err: %v", app, profile, err)
		return nil, err
	}
	return cfs, nil
}
func FindAppProfiles(app string) ([]string, error) {

	if len(app) == 0 {
		return nil, ErrInvalidParam
	}
	cf := db.Config{App: app}
	cfs, err := cf.ReadAppProfiles(app)

	if err != nil {
		log.Printf("find configs by {app: %s} failed, err: %v", app, err)
		return nil, err
	}
	return cfs, nil
}

func FindConfigById(id int64) (*db.Config, error) {
	if id <= 0 {
		return nil, ErrInvalidParam
	}

	cf := db.Config{Id: id}
	cfs, err := cf.ReadByInput()

	if err != nil {
		log.Printf("FindConfigById  {id: %d} failed, err: %v", id, err)
		return nil, err
	}

	if cfs == nil || len(cfs) == 0 {
		log.Printf("FindConfigById no matched config {id: %d} found.", id)
		return nil, nil
	}

	if len(cfs) > 1 {
		log.Printf("FindConfigById find config {id: %d} found too many matched record: %v", id, cfs)
	}

	return cfs[0], nil
}

func FindConfig(app string, profile string, key string) (*db.Config, error) {
	if len(app) == 0 || len(profile) == 0 || len(key) == 0 {
		return nil, ErrInvalidParam
	}

	cf := db.Config{App: app, Profile: profile, Key: key}
	cfs, err := cf.ReadByInput()

	if err != nil {
		log.Printf("find config {app: %s, profile: %s, key: %s} failed, err: %v", app, profile, key, err)
		return nil, err
	}

	if cfs == nil || len(cfs) == 0 {
		log.Printf("no matched config {app: %s, profile: %s, key: %s} found.", app, profile, key)
		return nil, nil
	}

	if len(cfs) > 1 {
		log.Printf("find config {app: %s, profile: %s, key: %s} found too many matched record: %v", app, profile, key, cfs)
	}

	return cfs[0], nil
}

func CountConfig(app string, profile string, key string) (int64, error) {
	if len(app) == 0 || len(profile) == 0 || len(key) == 0 {
		return 0, ErrInvalidParam
	}

	cf := db.Config{App: app, Profile: profile, Key: key}
	count, err := cf.Count()

	if err != nil {
		log.Printf("count config {app: %s, profile: %s, key: %s} failed, err: %v", app, profile, key, err)
		return 0, err
	}

	return count, nil
}

func CreateConfigByUser(uid int64, username string, app string, profile string, key string, value string, validator string, desc string) (*db.Config, error) {
	if uid <= 0 || len(username) == 0 || len(app) == 0 || len(profile) == 0 || len(key) == 0 {
		return nil, ErrInvalidParam
	}

	isManager := false
	if user.CheckManagePermission(uid, app) {
		isManager = true
	}

	return CreateConfig(username, app, profile, key, value, validator, desc, isManager)
}

func CreateConfig(username string, app string, profile string, key string, value string, validator string, desc string, isManager bool) (*db.Config, error) {
	if len(username) == 0 || len(app) == 0 || len(profile) == 0 || len(key) == 0 {
		return nil, ErrInvalidParam
	}

	count, err := CountConfig(app, profile, key)
	if err != nil {
		return nil, err
	}
	if count > 0 {
		return nil, ErrExists
	}

	currTime := time.Now()

	initVersion := Init_Version
	// Create config history
	cfhist := db.ConfigHistory{App: app, Profile: profile, Key: key}

	maxVersionHis, err := cfhist.ReadMaxVersionHistory()

	if err == nil && maxVersionHis != nil {
		initVersion = maxVersionHis.CurVersion + 1
	}

	cfhist.CurValue = value
	cfhist.CurVersion = initVersion
	cfhist.Modifier = username
	cfhist.ModifyTime = currTime
	cfhist.Validator = validator
	cfhist.Description = desc
	cfhist.OperateType = int8(db.OperationType_Create)

	err = cfhist.Create()
	if err != nil {
		log.Printf("create config history {%v} failed. err: %v", cfhist, err)
		return nil, err
	}

	// Create config
	cf := db.Config{App: app, Profile: profile, Key: key, Value: value, Validator: validator, Description: desc}
	cf.Version = initVersion
	cf.Modifier = username
	cf.ModifyTime = currTime
	cf.ApproveType = int8(db.ApproveType_NotApprove)

	if isManager {
		cf.Approver = username
		cf.ApproveType = int8(db.ApproveType_PASS)
	}

	err = cf.Create()
	if err != nil {
		log.Printf("create config {%v} failed. err: %v", cf, err)
		return nil, err
	}

	return &cf, nil
}

func DeleteConfig(user string, app string, profile string, key string) error {
	if len(user) == 0 || len(app) == 0 || len(profile) == 0 || len(key) == 0 {
		return ErrInvalidParam
	}

	dbcf, err := FindConfig(app, profile, key)
	if err != nil {
		return err
	}
	if dbcf == nil {
		return ErrNotExists
	}

	cfhist := db.ConfigHistory{App: dbcf.App, Profile: dbcf.Profile, Key: dbcf.Key}
	cfhist.PreVersion = dbcf.Version
	cfhist.PreValue = dbcf.Value
	cfhist.CurVersion = dbcf.Version + 1
	cfhist.Modifier = user
	cfhist.ModifyTime = time.Now()
	cfhist.Validator = dbcf.Validator
	cfhist.Description = dbcf.Description
	cfhist.OperateType = int8(db.OperationType_Delete)

	err = cfhist.Create()
	if err != nil {
		log.Printf("create config history {%v} failed. err: %v", cfhist, err)
		return err
	}

	num, err := dbcf.Delete()
	if err != nil {
		log.Printf("delete config {%v} failed. err: %v", dbcf, err)
		return err
	}

	if num != 1 {
		log.Printf("deleted config {%v} record count not as expected: %d", dbcf, num)
	}
	globals.GetEtcdClient().DeleteKeys(etcd.NewConfigKey(app, profile, key).FullPath)

	return nil
}

func UpdateConfigByUser(uid int64, username string, app string, profile string, key string, value string, validator string, desc string) (*db.Config, error) {
	if uid <= 0 || len(username) == 0 || len(app) == 0 || len(profile) == 0 || len(key) == 0 {
		return nil, ErrInvalidParam
	}

	isManager := false
	if user.CheckManagePermission(uid, app) {
		isManager = true
	}

	return UpdateConfig(username, app, profile, key, value, validator, desc, isManager)
}

func UpdateConfig(username string, app string, profile string, key string, value string, validator string, desc string, isManager bool) (*db.Config, error) {
	if len(username) == 0 || len(app) == 0 || len(profile) == 0 || len(key) == 0 {
		return nil, ErrInvalidParam
	}

	dbcf, err := FindConfig(app, profile, key)
	if err != nil {
		return nil, err
	}

	if dbcf == nil {
		return nil, ErrNotExists
	}

	version := dbcf.Version + 1
	currTime := time.Now()

	cfhist := db.ConfigHistory{App: dbcf.App, Profile: dbcf.Profile, Key: dbcf.Key}
	cfhist.CurValue = value
	cfhist.PreValue = dbcf.Value
	cfhist.CurVersion = version
	cfhist.PreVersion = dbcf.Version
	cfhist.Modifier = dbcf.Modifier
	cfhist.ModifyTime = currTime
	cfhist.Validator = validator
	cfhist.Description = desc
	cfhist.OperateType = int8(db.OperationType_Modify)

	err = cfhist.Create()
	if err != nil {
		log.Printf("create config history {%v} failed. err: %v", cfhist, err)
		return nil, err
	}

	// Create config
	dbcf.Value = value
	dbcf.Validator = validator
	dbcf.Description = desc
	dbcf.Version = version
	dbcf.Modifier = username
	dbcf.ModifyTime = currTime
	dbcf.Approver = ""
	dbcf.ApproveType = int8(db.ApproveType_NotApprove)

	if isManager {
		dbcf.Approver = username
		dbcf.ApproveType = int8(db.ApproveType_PASS)
	}

	num, err := dbcf.Update()
	if err != nil {
		log.Printf("update config {%v} failed. err: %v", dbcf, err)
		return nil, err
	}
	if num != 1 {
		log.Printf("updated config {%v} affect count not as expect: %d", dbcf, num)
	}

	return dbcf, nil
}

func RollbackConfigByUser(uid int64, username string, app string, profile string, key string, version int64) (*db.Config, error) {
	if uid <= 0 || len(username) == 0 || len(app) == 0 || len(profile) == 0 || len(key) == 0 || version <= 0 {
		return nil, ErrInvalidParam
	}

	isManager := false
	if user.CheckManagePermission(uid, app) {
		isManager = true
	}

	return RollbackConfig(username, app, profile, key, version, isManager)
}

func RollbackConfig(username string, app string, profile string, key string, version int64, isManager bool) (*db.Config, error) {
	if len(username) == 0 || len(app) == 0 || len(profile) == 0 || len(key) == 0 || version <= 0 {
		return nil, ErrInvalidParam
	}

	cfHist := db.ConfigHistory{App: app, Profile: profile, Key: key, CurVersion: version}
	dbCfHists, err := cfHist.ReadByInput()
	if err != nil {
		log.Printf("find config history {%v} failed. err: %v", cfHist, err)
		return nil, err
	}
	if dbCfHists == nil || len(dbCfHists) == 0 {
		log.Printf("can not found config history: %v", cfHist)
		return nil, ErrNotExists
	}
	if len(dbCfHists) > 1 {
		log.Printf("too many matched config history {%v}: %v", cfHist, dbCfHists)
		return nil, ErrNotExists
	}

	dbCfHist := dbCfHists[0]

	dbcf, err := FindConfig(app, profile, key)
	if err != nil {
		return nil, err
	}
	if dbcf == nil {
		return nil, ErrNotExists
	}

	newVersion := dbcf.Version + 1
	currTime := time.Now()
	desc := "Rollback from version " + strconv.FormatInt(version, 10)

	// Create config history
	cfhist := db.ConfigHistory{App: app, Profile: profile, Key: key}
	cfhist.CurValue = dbCfHist.CurValue
	cfhist.CurVersion = newVersion
	cfhist.PreValue = dbcf.Value
	cfhist.PreVersion = dbcf.Version
	cfhist.Modifier = username
	cfhist.ModifyTime = currTime
	cfhist.Validator = dbCfHist.Validator

	cfhist.Description = desc
	cfhist.OperateType = int8(db.OperationType_Rollback)

	err = cfhist.Create()
	if err != nil {
		log.Printf("create config history {%v} failed. err: %v", cfhist, err)
		return nil, err
	}

	// Create config
	cf := db.Config{App: app, Profile: profile, Key: key, Value: dbCfHist.CurValue, Validator: cfhist.Validator, Description: desc}
	cf.Version = newVersion
	cf.Modifier = username
	cf.ModifyTime = currTime

	cf.ApproveType = int8(db.ApproveType_NotApprove)
	if isManager {
		cf.Approver = username
		cf.ApproveType = int8(db.ApproveType_PASS)
	}

	err = cf.CreateOrUpdate()
	if err != nil {
		log.Printf("create config {%v} failed. err: %v", cf, err)
		return nil, err
	}

	return &cf, nil
}

func ApproveConfigByUser(uid int64, username string, app string, profile string, key string, version int64) error {
	if uid <= 0 || len(username) == 0 || len(app) == 0 || len(profile) == 0 || len(key) == 0 || version <= 0 {
		return ErrInvalidParam
	}

	if !user.CheckManagePermission(uid, app) {
		log.Printf("approve config {app: %s, profile: %s, key: %s, version: %d} failed due to permission denied: %d", app, profile, key, version, username)
		return ErrPermissionDenied
	}

	dbcf, err := FindConfig(app, profile, key)
	if err != nil {
		return err
	}
	if dbcf == nil {
		log.Printf("no matched config {app: %s, profile: %s, key: %s, version: %d} to approve: %s", app, profile, key, version, username)
		return ErrNotExists
	}
	if version != dbcf.Version {
		log.Printf("approve version is old {%v}. approve version: %d, user: %s", dbcf, version, username)
		return ErrVersionOld
	}

	dbcf.Approver = username
	dbcf.ApproveType = int8(db.ApproveType_PASS)

	num, err := dbcf.Update()
	if err != nil {
		log.Printf("update config {%v} to approved failed. err: %v", dbcf, err)
		return err
	}

	if num != 1 {
		log.Printf("updated approve config {%v} record count not as expected: %d", dbcf, num)
	}

	return nil
}

func QueryConfigHistory(app string, profile string, key string) ([]*db.ConfigHistory, error) {
	if len(app) == 0 || len(profile) == 0 || len(key) == 0 {
		return nil, ErrInvalidParam
	}

	cfhist := db.ConfigHistory{App: app, Profile: profile, Key: key}

	cfhistList, err := cfhist.ReadByInput()
	if err != nil {
		log.Printf("QueryConfigHistory config history {%v} failed. err: %v", cfhist, err)
		return nil, err
	}
	return cfhistList, nil
}

func DeleteConfigByApp(app string) error {
	config := db.Config{App: app}
	config.DeleteByInput()

	configHistory := db.ConfigHistory{App: app}
	configHistory.DeleteByInput()

	globals.GetEtcdClient().DeleteKeys(etcd.KeySeprator + app + etcd.KeySeprator)
	return nil
}

func DeleteConfigByAppProfile(app, profile string) error {
	config := db.Config{App: app, Profile: profile}
	config.DeleteByInput()

	configHistory := db.ConfigHistory{App: app, Profile: profile}
	configHistory.DeleteByInput()
	globals.GetEtcdClient().DeleteKeys(etcd.GetKeyPrefix(app, profile, etcd.KeyType_Profile))
	return nil
}
func CopyConfig(app, srcProfile, destProfile, creator string, creatorUid int64) error {
	isManager := false
	if user.CheckManagePermission(creatorUid, app) {
		isManager = true
	}
	config := db.Config{App: app, Profile: srcProfile}
	return config.CopyTo(destProfile, creator, isManager)
}
