package server

import (
	"errors"
	"fmt"
	"log"
	"net"
	"strconv"

	"github.com/yy-java/cnt2/grpc/pb"
	"github.com/yy-java/cnt2/service/config"
	"github.com/yy-java/cnt2/service/publish"
	"github.com/yy-java/cnt2/service/register"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// ConfigCenterServiceServer desc
type ConfigCenterServiceServer struct {
	ServerName string
	Port       int
}

// Register client and return nodeId
func (s *ConfigCenterServiceServer) RegisterClient(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	var (
		pid    int   = 0
		err    error = nil
		nodeId int64 = 0
	)

	resp := &pb.RegisterResponse{Result: 0, NodeId: ""}

	if len(req.Pid) > 0 {
		pid, err = strconv.Atoi(req.Pid)
		if err != nil {
			pid = 0
		}
	}

	nodeId, err = register.RegisterClient(req.App, req.Profile, req.ServerIp, pid)
	if nodeId > 0 {
		resp.Result = 1
		resp.NodeId = fmt.Sprintf("%d", nodeId)
	} else {
		if err != nil {
			resp.NodeId = ""
			if err == register.AppNotExistErr {
				resp.Result = 404
			} else if err == register.ParamErr {
				resp.Result = 400
			} else {
				resp.Result = 500
			}
		} else {
			resp.Result = 500
		}
	}

	log.Printf("RegisterClient [%v] result=[%v]\n", req, resp)

	return resp, nil
}

// QueryAll
func (s *ConfigCenterServiceServer) QueryAll(ctx context.Context, req *pb.QueryRequest) (*pb.QueryResponse, error) {
	log.Printf("QueryAll %v\n", req)

	configs, err := config.FindAllConfig(req.App, req.Profile)
	if err != nil {
		return nil, err
	}

	var respMsgs []*pb.ResponseMessage
	for i, cf := range configs {
		if cf.PublishedVersion > 0 {
			respMsgs = append(respMsgs, &pb.ResponseMessage{Key: configs[i].Key, Path: "/test", Profile: cf.Profile,
				Value: cf.PublishedValue, Version: cf.PublishedVersion})
		}
	}

	resp := &pb.QueryResponse{Result: respMsgs}

	return resp, nil
}

// QueryKey
func (s *ConfigCenterServiceServer) QueryKey(ctx context.Context, req *pb.QueryRequest) (*pb.ResponseMessage, error) {
	log.Printf("QueryKey %v\n", req)

	configHistory := config.FindConfigHistoryByKeyAndVersion(req.App, req.Profile, req.Key, req.KeyVersion)

	// 返回path
	var respMsgs *pb.ResponseMessage
	if configHistory != nil {
		respMsgs = &pb.ResponseMessage{Key: configHistory.Key, Path: "/test", Profile: configHistory.Profile,
			Value: configHistory.CurValue, Version: configHistory.CurVersion}
	} else { // 返回nil 还是 默认值？？
		respMsgs = &pb.ResponseMessage{}
	}
	resp := respMsgs

	return resp, nil
}

// ValueChangeResultNotify
func (s *ConfigCenterServiceServer) ValueChangeResultNotify(ctx context.Context, req *pb.ValueChangeResultRequest) (*pb.ValueChangeResultResponse, error) {
	resp := &pb.ValueChangeResultResponse{}

	nodeId, err := strconv.Atoi(req.NodeId)
	deployId, err := strconv.Atoi(req.DeployId)
	ret, err := publish.ValueChangeResultNotify(int64(nodeId), req.App, req.Profile, req.Key, int64(deployId), req.Version, req.Result)

	if err != nil {
		if err == publish.ParamErr {
			resp.Status = 400
			resp.Msg = ""
		} else if err == publish.NotFound {
			resp.Status = 404
			resp.Msg = ""
		} else {
			resp.Status = 500
			resp.Msg = ""
		}
	} else if ret == 1 {
		resp.Status = 1
		resp.Msg = ""
	} else {
		resp.Status = 0
		resp.Msg = ""
	}

	log.Printf("ValueChangeResultNotify [%v] result=[%v]\n", req, resp)

	return resp, nil
}

func NewServer(port int) *ConfigCenterServiceServer {
	s := new(ConfigCenterServiceServer)
	s.ServerName = "ConfigCenterServer"
	s.Port = port
	return s
}

func Start(server *ConfigCenterServiceServer) error {
	port := ":0"
	if server.Port > 0 {
		port = fmt.Sprintf(":%d", server.Port)
	} else {
		return errors.New("error server port")
	}

	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("gRPC server failed to listen: %v", err)
		return err
	}
	s := grpc.NewServer()

	pb.RegisterConfigCenterServiceServer(s, server)

	// Register reflection service on gRPC server.
	reflection.Register(s)

	if err := s.Serve(lis); err != nil {
		log.Fatalf("gRPC server failed to serve: %v", err)
		return err
	}
	return nil
}
