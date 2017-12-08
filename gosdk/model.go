package gosdk

import (
	"github.com/boltdb/bolt"
	"time"
)

const COMMON = "common"

type StringArray []string

type Cnt2Service struct {
	*ClientConfig
}

type GrpcServerInfo struct {
	ServerIP map[int]string `json:"serverIP"`
	Port     int            `json:"port"`
	GroupId  int            `json:"groupId"`
}

type ClientConfig struct {
	Endpoints     []string
	DialTimeout   time.Duration
	LocalFilePath string //指定本地缓存的存储地址，便于在配置中心不可用时，可以使用之前的配置
	App           string
	Profile       string
	EnableCommon  bool
}

type LocalStore struct {
	Db *bolt.DB
}
type ConfigPublishInfo struct {
	PublishId    int64       `json:"publishId"`
	Key          string      `json:"key"`
	Version      int64       `json:"version"`
	PublishNodes StringArray `json:"publishNodes"`
}
type Config struct {
	App         string            `json:"app"`
	Profile     string            `json:"profile"`
	Key         string            `json:"key"`
	Version     int64             `json:"version"`
	PublishInfo ConfigPublishInfo `json:"publishInfo"`
	Value       string            `json:"value"`
}
type Node struct {
	NodeId       string `json:"nodeId"`
	App          string `json:"app"`
	Profile      string `json:"profile"`
	Sip          string `json:"sip"`
	Pid          int    `json:"pid"`
	RegisterTime int64  `json:"registerTime"`
}
