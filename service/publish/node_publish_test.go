package publish

import (
	"testing"
	"github.com/yy-java/cnt2/service/register"
	"github.com/yy-java/cnt2/etcd"
	"github.com/yy-java/cnt2/service/models"
	"github.com/astaxie/beego/orm"
	"time"
	"strconv"
	"log"
	"github.com/yy-java/cnt2/db"
)

const (
	DB_Test_Name       = "default"
	DB_TEST_DRIVER     = "mysql"
	DB_TEST_DataSource = "cnt2_db_user:q0NUVMca1@tcp(58.215.143.133:6307)/cnt2_db?charset=utf8"
	DB_TEST_MaxIdle    = 30
)

func init() {
	// set default database
	orm.RegisterDataBase(DB_Test_Name, DB_TEST_DRIVER, DB_TEST_DataSource, DB_TEST_MaxIdle)
	orm.Debug = true
}

var (
	testEndpoints = []string{"172.27.141.11:2379", "172.27.141.12:2379", "172.27.141.13:2379"}
	defaultConfig = etcd.EtcdClientConfig{Endpoints: testEndpoints, DialTimeout: 3 * time.Second, RequestTimeout: 3 * time.Second}
	etcdClient, _ = etcd.NewEtcdClient(defaultConfig)

	app              = "app1"
	profile          = "wx"
	nodeIdList       = []int64{1111, 2222, 3333}
	pidList          = []int{11, 22, 33}
	ipList           = []string{"127.0.0.1", "127.0.0.2", "127.0.0.3"}
	key              = "somePropertyKey"
	publishId  int64 = 199
	version    int64 = 1
)

func registerEtcdNode() {
	register.ClientOnline(etcdClient, app, profile, nodeIdList[0])
	register.ClientOnline(etcdClient, app, profile, nodeIdList[1])
	register.ClientOnline(etcdClient, app, profile, nodeIdList[2])

	ret := ListAllNode(etcdClient, app, profile)
	log.Printf("node:%v", ret)
}
func registerDbNode() {
	register.RegisterClient(app, profile, ipList[0], pidList[0])
	register.RegisterClient(app, profile, ipList[1], pidList[1])
}
func prepareNode() {
	//registerDbNode()
	registerEtcdNode()
}
func deleteDbNodePublish() {
	for _, nodeId := range nodeIdList {
		node := db.Node{Id: nodeId}
		node.Delete()
	}
}
func TestPartitionPublishNode(t *testing.T) {
	prepareNode()
	publishNodeIdList := []string{strconv.FormatInt(nodeIdList[0], 10)}
	PartitionPublishNode(etcdClient, app, profile, key, publishId, version, publishNodeIdList)

	var configNode models.ConfigNode
	getConfigValue(etcdClient, app, profile, key, &configNode)
	t.Logf("etcd config publish value:%v", configNode)
	deleteDbNodePublish()
}
func TestAllPublishNode(t *testing.T) {
	prepareNode()
	var configNode models.ConfigNode
	getConfigValue(etcdClient, app, profile, key, &configNode)
	t.Logf("etcd config publish value:%v", configNode)

	AllPublishNode(etcdClient, app, profile, key, publishId, version)
	getConfigValue(etcdClient, app, profile, key, &configNode)
	t.Logf("etcd config publish value:%v", configNode)
	deleteDbNodePublish()
}
