package node

import (
	"log"
	"github.com/yy-java/cnt2/db"
	"github.com/yy-java/cnt2/service/errors"
	"github.com/yy-java/cnt2/service/user"
)

func FindByAppAndProfile(uid int64, app, profile string) ([]*db.Node, error) {

	log.Printf("FindByAppAndProfile app:%s, profile:%s", app, profile)

	if len(app) == 0 || len(profile) == 0 {
		return nil, errors.ErrInvalidParam
	}
	permission := user.CheckPermission(uid, app)
	if permission <= 0 {
		return nil, errors.ErrPermissionDenied
	}

	nodes, err := (&db.Node{}).FindAllNode(app, profile)
	if err != nil {
		log.Printf("FindByAppAndProfile  app:%s, profile:%s failed: %v", app, profile, err)
		return nil, errors.ErrServerErr
	}
	return nodes, nil
}

/**
返回已经发布的节点ID（node_publish的node_id）
**/
func FindPublishedNode(app, profile, key string, version int64) ([]string, error) {

	log.Printf("FindPublishedNode app:%s, profile:%s, key:%s,version:%s", app, profile, key, version)

	if len(app) == 0 || len(profile) == 0 || len(key) == 0 || version < 0 {
		return nil, errors.ErrInvalidParam
	}

	nodes, err := (&db.NodePublish{}).FindPublishedNode(app, profile, key, version)
	if err != nil {
		log.Printf("FindByAppAndProfile  app:%s, profile:%s failed: %v", app, profile, err)
		return nil, errors.ErrServerErr
	}
	return nodes, nil
}
