package gosdk

import (
	"encoding/json"
	"fmt"
	"github.com/coreos/etcd/clientv3"
	"golang.org/x/net/context"
	"log"
	"strings"
)

const (
	key_seprator         = "/"
	grpcServer_key_preix = key_seprator + "grpcservers"
)

func GenWatchCnt2Key(app, profile string) string {
	return fmt.Sprintf("/%s/profiles/%s/", app, profile)
}
func GenWatchCnt2CommonKey(app string) string {
	return fmt.Sprintf("/%s/profiles/common/", app)
}
func GenWatchGrpcServerKey() string {
	return grpcServer_key_preix
}
func GenRegisterNodeKey(app, profile, nodeId string) string {
	return fmt.Sprintf("/%s/nodes/%s/%s", app, profile, nodeId)
}

func (config *Config) FromJson(val string) *Config {
	err := json.Unmarshal([]byte(val), &config)
	if err != nil {
		return nil
	}
	return config
}
func (config *Config) ToJson() string {
	val, err := json.Marshal(config)
	if err != nil {
		return ""
	}
	return string(val)
}

/**
return app ,profile,key
*/
func ExtractCntInfo(nodeKey string) (string, string, string, bool) {
	result := strings.Split(nodeKey, key_seprator)
	if len(result) < 4 {
		return "", "", "", false
	}
	return result[1], result[3], result[4], true
}
func (array StringArray) Remove(target string) StringArray {
	for i, v := range array {
		if v == target {
			array = append(array[:i], array[i+1:]...)
			return array
		}
	}
	return array
}
func (array StringArray) Contains(val string) bool {
	if array == nil || len(array) <= 0 {
		return false
	}
	for _, i := range []string(array) {
		if i == val {
			return true
		}
	}
	return false
}
func ParseGrpcServerInfoFromEtcd(client *clientv3.Client) ([]GrpcServerInfo, error) {
	resp, err := client.Get(context.Background(), GenWatchGrpcServerKey(), clientv3.WithPrefix())
	if err != nil {
		return nil, err
	}
	var grpcServerInfos []GrpcServerInfo
	for _, kv := range resp.Kvs {
		gi := GrpcServerInfo{}
		if err := json.Unmarshal(kv.Value, &gi); err == nil {
			grpcServerInfos = append(grpcServerInfos, gi)
		}
	}
	log.Printf("grpcServerInfo:%v", grpcServerInfos)
	return grpcServerInfos, nil
}
