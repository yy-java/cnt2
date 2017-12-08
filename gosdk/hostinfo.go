package gosdk

/*
通过运维脚本获取服务器的host信息，
获取本机IP信息
*/

import (
	"bufio"
	//"io"
	"net"
	"os"
	"strconv"
	"strings"
)

const (
	hostInfoPath = "/home/dspeak/yyms/hostinfo.ini"
)

var (
	NetNameMap = make(map[string]int)
	NetTypeMap = make(map[int]string)
)

type NetType struct {
	Name string
	Val  int
}

type HostIP struct {
	Type NetType
	Ip   string
}

type HostInfo struct {
	Ips     []HostIP
	GroupId int
}

func InitHostInfo() HostInfo {
	/**
	 * 电信
	 */
	NetNameMap["CTL"] = 0x1
	/**
	 * 网通
	 */
	NetNameMap["CNC"] = 0x2
	/**
	 * 铁通
	 */
	NetNameMap["CNII"] = 0x4
	/**
	 * 教育
	 */
	NetNameMap["EDU"] = 0x8
	/**
	 * 长城
	 */
	NetNameMap["WBN"] = 0x16
	/**
	 * 移动网
	 */
	NetNameMap["MOB"] = 0x32
	NetNameMap["BGP"] = 0x64
	NetNameMap["HK"] = 0x128
	NetNameMap["BRA"] = 0x257

	NetNameMap["INTRANET"] = 0x999

	for k, v := range NetNameMap {
		NetTypeMap[v] = k
	}
	return getHostInfo(hostInfoPath)
}

func getHostInfo(path string) HostInfo {
	var hi HostInfo
	if len(path) == 0 {
		path = hostInfoPath
	}
	properties := propertiesReader(path)

	if len(properties) == 0 {
		return hi
	}

	iplist, ok := properties["ip_isp_list"]
	if !ok {
		return hi
	}

	iplists := strings.Split(iplist, ",")
	for _, ipw := range iplists {
		s := strings.Split(ipw, ":")
		if len(s) == 2 {
			ip := s[0]
			nt := s[1]
			ntVal, ok := NetNameMap[nt]

			if !ok {
				ntVal = 0
				continue
			}
			netType := NetType{Name: nt, Val: ntVal}
			host := HostIP{netType, ip}
			hi.Ips = append(hi.Ips, host)
		}
	}
	grpId, ok := properties["pri_group_id"]
	if !ok {
		grpId = "0"
	}
	hi.GroupId, _ = strconv.Atoi(grpId)
	return hi
}

func propertiesReader(filePath string) map[string]string {
	ret := make(map[string]string)
	f, err := os.Open(filePath)
	if err != nil {
		return ret //fileNotExists
	}
	defer f.Close()
	r := bufio.NewReader(f)
	for {
		line, _, err := r.ReadLine()
		if err != nil {
			/*if err == io.EOF {
				break
			}*/
			break
		}

		s := strings.TrimSpace(string(line))
		if strings.Index(s, "#") == 0 || len(s) == 0 {
			continue
		}
		n1 := strings.Index(s, "=")
		if n1 == len(s)-1 {
			continue
		}
		key := strings.TrimSpace(s[0:n1])
		value := strings.TrimSpace(s[n1+1:])
		ret[key] = value
	}
	return ret
}

func GetMyIPInfo() []string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return make([]string, 0)
	}
	ret := make([]string, 0, len(addrs))
	for _, address := range addrs {
		// 检查ip地址判断是否回环地址
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				ip := ipnet.IP.To4().String()
				ip = strings.TrimSpace(ip)
				if len(ip) > 0 {
					ret = append(ret, ip)
				}
			}
		}
	}
	return ret
}
