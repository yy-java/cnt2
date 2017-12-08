package gosdk

import (
	"google.golang.org/grpc/naming"
	"log"
	"strconv"
	"strings"
)

var grpcServerChan = GrpcServerChan{ch: make(chan []GrpcServerInfo)}

type GrpcServerChan struct {
	ch          chan []GrpcServerInfo
	bestAddress StringArray
}

type GrpcServerResolver struct {
	serverInfo []GrpcServerInfo
}

type GrpcServerWatcher struct {
	*GrpcServerResolver
	isInited    bool
	initAddress StringArray
}

func (re *GrpcServerResolver) Resolve(target string) (naming.Watcher, error) {
	return &GrpcServerWatcher{re, false, strings.Split(target, ";")}, nil
}

func (w *GrpcServerWatcher) Next() ([]*naming.Update, error) {
	if w.isInited == false {
		w.isInited = true

		var namingUpdate []*naming.Update
		for _, address := range w.initAddress {
			namingUpdate = append(namingUpdate, &naming.Update{Op: naming.Add, Addr: address})
			grpcServerChan.bestAddress = append(grpcServerChan.bestAddress, address)
		}

		return namingUpdate, nil

	}
	for serverInfo := range grpcServerChan.ch {
		log.Printf("the New grpcServer: %s", serverInfo)
		if namingUpdate, ok := buildNamingUpdate(serverInfo); ok {
			return namingUpdate, nil
		}
	}
	return nil, nil
}
func buildNamingUpdate(serverInfo []GrpcServerInfo) ([]*naming.Update, bool) {
	if newBestAddress, ok := ChooseBestAddress(serverInfo); ok {
		log.Printf("the best address: %s", newBestAddress)
		var updates []*naming.Update

		if len(grpcServerChan.bestAddress) <= 0 {
			for _, address := range newBestAddress {
				updates = append(updates, &naming.Update{Op: naming.Add, Addr: address})
				grpcServerChan.bestAddress = append(grpcServerChan.bestAddress, address)
			}

		} else {

			for _, newAddress := range newBestAddress {
				var find bool
				for _, oldAddress := range grpcServerChan.bestAddress {
					if oldAddress == newAddress {
						find = true
						break
					}
				}
				if !find {
					updates = append(updates, &naming.Update{Op: naming.Add, Addr: newAddress})
					grpcServerChan.bestAddress = append(grpcServerChan.bestAddress, newAddress)
				}
			}
			for _, oldAddress := range grpcServerChan.bestAddress {
				var find bool
				for _, newAddress := range newBestAddress {
					if oldAddress == newAddress {
						find = true
						break
					}
				}
				if !find {
					updates = append(updates, &naming.Update{Op: naming.Delete, Addr: oldAddress})
					grpcServerChan.bestAddress = grpcServerChan.bestAddress.Remove(oldAddress)
				}
			}
		}
		if len(updates) > 0 {
			return updates, true
		}
	}
	return nil, false
}

func (re *GrpcServerWatcher) Close() {
	//nothing to do
}

//FIXME can be better?
func ChooseBestAddress(grpcServerInfos []GrpcServerInfo) ([]string, bool) {

	if len(hostInfo.Ips) == 0 {
		for k := range NetTypeMap {
			//不能是内网
			if k == 2457 {
				continue
			}
			if v, ok := grpcServerInfos[0].ServerIP[k]; ok {
				return StringArray{v + ":" + strconv.Itoa(grpcServerInfos[0].Port)}, true
			}
		}
	}
	var sameGroup, bestChoice, sameNetType, random []string
	for _, serverInfo := range grpcServerInfos {
		//同机房
		if hostInfo.GroupId == serverInfo.GroupId {
			for _, ip := range hostInfo.Ips {
				if v, ok := serverInfo.ServerIP[ip.Type.Val]; ok {
					//优先同运营商
					bestChoice = append(bestChoice, v+":"+strconv.Itoa(serverInfo.Port))
					break
				}
			}
			if len(bestChoice) == 0 {
				for k := range NetTypeMap {
					//不能是内网
					if k == 2457 {
						continue
					}
					if v, ok := serverInfo.ServerIP[k]; ok {
						sameGroup = append(sameGroup, v+":"+strconv.Itoa(serverInfo.Port))
						break
					}
				}
			}
			//不同机房
		} else {
			for _, ip := range hostInfo.Ips {
				if v, ok := serverInfo.ServerIP[ip.Type.Val]; ok {
					//运营商
					sameNetType = append(sameNetType, v+":"+strconv.Itoa(serverInfo.Port))
					break
				}
				if len(sameNetType) == 0 {
					//随机
					for k := range NetTypeMap {
						//不能是内网
						if k == 2457 {
							continue
						}
						if v, ok := serverInfo.ServerIP[k]; ok {
							random = append(random, v+":"+strconv.Itoa(serverInfo.Port))
						}
					}
				}
			}
		}
	}
	if len(bestChoice) > 0 {
		return bestChoice, true
	}
	if len(sameGroup) > 0 {
		return sameGroup, true
	}
	if len(sameNetType) > 0 {
		return sameNetType, true
	}
	if len(random) > 0 {
		return random, true
	}
	return random, false
}
