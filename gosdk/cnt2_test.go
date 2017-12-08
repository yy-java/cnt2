package gosdk

import (
	"fmt"
	//	"github.com/coreos/etcd/clientv3"
	"testing"
	"time"
)

type TestListenter struct{}

func (t *TestListenter) HandlePutEvent(config *Config) error {
	fmt.Printf("put key: %s ; newValue: %s; version:%s \n", config.Key, config.Value, config.Version)
	return nil
}
func (t *TestListenter) HandleDeleteEvent(config *Config) error {
	fmt.Printf("delete app:%s, profile:%s, key: %s \n", config.App, config.Profile, config.Key)
	return nil
}
func Test_NewClient(t *testing.T) {
	clientConfig := ClientConfig{
		Endpoints:     []string{"61.147.187.152:2379", "61.147.187.142:2379", "61.147.187.150:2379"},
		App:           "demo",
		Profile:       "development",
		LocalFilePath: "",
		EnableCommon:  false}
	cnt2Service, err := Start(&clientConfig)
	if err != nil {
		fmt.Printf(" error:%v \n", err)
	}
	if cnt2Service != nil {
		cnt2Service.RegisterListener(&TestListenter{}, "my_config_1", "test_config", "common_config")

		fmt.Printf("client:%v, error:%v \n", cnt2Service, err)
		time.Sleep(time.Minute * 10)

		cnt2Service.Close()
	}
}

func Test_GetConfig(t *testing.T) {
	for i := 0; i < 5; i++ {
		go get(i)
	}
	time.Sleep(time.Second * 2)
}
func get(i int) {
	cnt2Service := &Cnt2Service{&ClientConfig{App: "demo", Profile: "development"}}
	val, _ := cnt2Service.GetConfig("my_config_1")
	fmt.Println(i)
	fmt.Println("my_config_1:" + val)
}
