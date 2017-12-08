package cnt2sdk

import (
	"fmt"
	"github.com/coreos/etcd/mvcc/mvccpb"
	"testing"
	"time"
	"github.com/yy-java/cnt2/etcd"
)

var (
	//	testEndpoints = []string{"221.228.91.155:2379", "221.228.83.119:2379", "113.108.65.32:2379"}
	testEndpoints = []string{"61.147.187.152:2379", "61.147.187.142:2379", "61.147.187.150:2379"}
)

type TestListenter struct{}

func (t *TestListenter) Callback(key, newValue, oldValue string, event mvccpb.Event_EventType) error {
	fmt.Println("change key: " + key + " newValue: " + newValue)
	return nil
}

func Test_GetConfig(t *testing.T) {
	config := XClientConfig{
		EtcdClientConfig: etcd.EtcdClientConfig{
			Endpoints: testEndpoints,
		},
		App:     "demo",
		Profile: "development",
	}
	client, err := NewXClient(config)
	if err != nil {
		t.Error(err)
	}
	if client == nil {
		return
	}
	k, ok := client.GetConfig("tks")
	fmt.Println(k, ok)
}
func Test_Listener(t *testing.T) {
	config := XClientConfig{
		EtcdClientConfig: etcd.EtcdClientConfig{
			Endpoints: testEndpoints,
		},
		App:     "demo",
		Profile: "development",
	}
	client, err := NewXClient(config)
	if err != nil {
		t.Error(err)
	}
	if client == nil {
		return
	}
	client.AddListener("s", &TestListenter{})
	time.Sleep(10000 * time.Second)

}
