package publish

import (
	"errors"
	"log"

	"github.com/yy-java/cnt2/db"
	"github.com/yy-java/cnt2/grpc/pb"
)

var (
	ParamErr  = errors.New("param error")
	NotFound  = errors.New("not found")
	StatusErr = errors.New("status error")
)

// 更新配置发布结果状态，失败返回0，成功返回1
func ValueChangeResultNotify(nodeId int64, app string, profile string, key string, deployId int64, version int64,
	result pb.ValueChangeResultRequest_ValueChangeResult) (int32, error) {
	if len(app) == 0 || len(profile) == 0 || len(key) == 0 {
		return 0, ParamErr
	}

	nodePublish := &db.NodePublish{NodeId: nodeId, App: app, Profile: profile, Key: key, PublishResult: -1, Version: version, Id: deployId}
	nodePublishs, err := nodePublish.ReadByInput()
	if err != nil {
		log.Printf("find nodePublish {%v} failed: %v", nodePublish, err)
		return 0, err
	}

	if len(nodePublishs) == 0 {
		return 0, NotFound
	}

	nodePublish = nodePublishs[0]
	if nodePublish.PublishResult != int8(db.PublishResult_Create) {
		return 0, StatusErr
	}

	var publishResult db.PublishResult = db.PublishResult_Create
	switch result {
	case pb.ValueChangeResultRequest_SUCCESS:
		publishResult = db.PublishResult_Success
	case pb.ValueChangeResultRequest_FAILED:
		publishResult = db.PublishResult_Fail
	default:
		publishResult = db.PublishResult_Fail
	}
	nodePublish.PublishResult = int8(publishResult)
	_, err = nodePublish.UpdatePublishResult()
	if err != nil {
		log.Printf("update nodePublish {%v} failed: %v", nodePublish, err)
		return 0, err
	}

	return 1, nil
}
