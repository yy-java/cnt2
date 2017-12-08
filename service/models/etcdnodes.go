package models

import "encoding/json"

type JsonObject interface {
	ToJson() (string, error)
	FromJson(string) error
}

func toJson(j JsonObject) (string, error) {
	jbuf, err := json.Marshal(j)
	if err != nil {
		return "", err
	}
	return string(jbuf), nil
}

func fromJson(j JsonObject, buf string) error {
	err := json.Unmarshal([]byte(buf), j)
	if err != nil {
		return err
	}
	return nil
}

// RegisterNode clientSdk Register in etcd
type RegisterNode struct {
	NodeId       string `json:"nodeId"`
	App          string `json:"app"`
	Profile      string `json:"profile"`
	Pid          int    `json:"pid"`
	Sip          string `json:"sip"`
	RegisterTime int64  `json:"registerTime"`
	ServerIP     string `json:"serverIp"`
}

func (g *RegisterNode) ToJson() (string, error) {
	return toJson(g)
}

func (g *RegisterNode) FromJson(jbuf string) error {
	return fromJson(g, jbuf)
}

type ConfigPublishInfo struct {
	PublishId    int64    `json:"publishId"`
	Key          string   `json:"key"`
	Version      int64    `json:"version"`
	PublishNodes []string `json:"publishNodes"`
}

type ConfigNode struct {
	App         string            `json:"app"`
	Profile     string            `json:"profile"`
	Key         string            `json:"key"`
	Version     int64             `json:"version"`
	PublishInfo ConfigPublishInfo `json:"publishInfo"`
}

func (g *ConfigNode) ToJson() (string, error) {
	return toJson(g)
}

func (g *ConfigNode) FromJson(jbuf string) error {
	return fromJson(g, jbuf)
}

// GrpcServerRegisterInfo grpc server info
type GrpcServerRegisterInfo struct {
	ServerIP map[int]string `json:"serverIP"`
	Port     int            `json:"port"`
	GroupId  int            `json:"groupId"`
}

func (g *GrpcServerRegisterInfo) ToJson() (string, error) {
	return toJson(g)
}

func (g *GrpcServerRegisterInfo) FromJson(jbuf string) error {
	return fromJson(g, jbuf)
}
