package config

import (
	"log"

	"github.com/yy-java/cnt2/db"
)

func FindConfigHistoryByKeyAndVersion(app string, profile string, key string, keyVersion int64) *db.ConfigHistory {
	
	log.Printf("FindConfigHistoryByKeyAndVersion app:%s, profile:%s, key:%s, keyVersion:%d\n", app, profile, key, keyVersion)
	
	if len(app) == 0 || len(profile) == 0 {
		return nil
	}
	
	configHistory := &db.ConfigHistory{App:app, Profile:profile, Key:key, CurVersion:keyVersion}
	configHistorys, err := configHistory.ReadByInput()
	if err != nil {
		log.Printf("find configHistory {%v} failed: %v", configHistory, err)
		return nil
	}
	
	if len(configHistorys) > 0 {
		return configHistorys[0]
	}
	return nil
}
