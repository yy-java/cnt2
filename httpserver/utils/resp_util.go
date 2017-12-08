package utils

import (
	"github.com/yy-java/cnt2/service/codes"
	"github.com/yy-java/cnt2/service/errors"
)

func Resp(data interface{}, err error) map[string]interface{} {
	resp := make(map[string]interface{})
	if nil != err {

		if newErr, ok := err.(*errors.ServiceError); ok {
			resp["code"] = newErr.Code()
			resp["msg"] = newErr.Error()
		} else {
			resp["code"] = codes.ServerError
			resp["msg"] = newErr.Error()
		}

		return resp
	}

	resp["code"] = codes.OK
	if data != nil {
		resp["data"] = data
	}
	return resp
}
