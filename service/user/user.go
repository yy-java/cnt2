package user

import (
	"log"

	. "github.com/yy-java/cnt2/db"
)

func SaveUserAuth(auth *UserAuth) error {
	if auth.Id > 0 {
		num, err := auth.Update()
		if err != nil {
			log.Printf("update UserAuth {%v} failed: %v", auth, err)
			return err
		}

		if num <= 0 {
			log.Printf("updated UserAuth {%v} and no record affect: %d", auth, num)
		}

		return nil
	}

	err := auth.Create()
	if err != nil {
		log.Printf("Save UserAuth {%v} error: %v", auth, err)
		return err
	}

	return nil
}

func FindUserAuthById(id int64) (*UserAuth, error) {
	auth := UserAuth{Id: id}
	err := auth.Read()
	if err != nil {
		log.Printf("find UserAuth {id: %d} failed: %v", id, err)
		return nil, err
	}

	return &auth, nil
}

func FindUserAuthByInput(auth *UserAuth) ([]*UserAuth, error) {
	uAuths, err := auth.ReadByInput()
	if err != nil {
		log.Printf("find UserAuth {%v} failed: %v", auth, err)
		return nil, err
	}

	return uAuths, nil
}

func RemoveUserAuthById(id int64) error {
	auth := UserAuth{Id: id}
	num, err := auth.Delete()
	if err != nil {
		log.Printf("delete UserAuth {id: %d} failed: %v", id, err)
		return err
	}

	if num <= 0 {
		log.Printf("deleted UserAuth {id: %d} and no record affect: %d", id, num)
	}

	return nil
}

func RemoveUserAuthByUidAndApp(uid int64, app string) error {
	auth := UserAuth{Uid: uid, App: app}
	num, err := auth.DeleteByInput()
	if err != nil {
		log.Printf("delete UserAuth {uid: %d, app: %s} failed: %v", uid, app, err)
		return err
	}

	if num <= 0 {
		log.Printf("deleted UserAuth {uid: %d, app: %s} and no record affect: %d", uid, app, num)
	}

	return nil
}
func CheckPermission(uid int64, app string) int8 {
	//查看是否是超级管理员
	auths, err := FindUserAuthByInput(&UserAuth{Uid: uid, Permission: int8(PermissionType_Admin)})
	if err == nil && len(auths) > 0 {
		return int8(PermissionType_Admin)
	}

	uAuths, err := FindUserAuthByInput(&UserAuth{Uid: uid, App: app})
	if err != nil {
		return 0
	}
	if uAuths == nil || len(uAuths) == 0 {
		return 0
	}

	return uAuths[0].Permission
}
func CheckManagePermission(uid int64, app string) bool {
	//查看是否是超级管理员
	auths, err := FindUserAuthByInput(&UserAuth{Uid: uid, Permission: int8(PermissionType_Admin)})
	if err == nil && len(auths) > 0 {
		return true
	}

	uAuths, err := FindUserAuthByInput(&UserAuth{Uid: uid, App: app})

	if err != nil {
		return false
	}
	if uAuths == nil || len(uAuths) == 0 {
		return false
	}

	return IsManager(uAuths[0].Permission)
}

func IsManager(permission int8) bool {
	switch permission {
	case int8(PermissionType_Dev):
		return false
	case int8(PermissionType_Admin):
		return true
	case int8(PermissionType_Manage):
		return true
	default:
		return false
	}

}
