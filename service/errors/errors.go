package errors

import "github.com/yy-java/cnt2/service/codes"

var (
	ErrServerErr        = New(codes.ServerError, "server not available")
	ErrInvalidParam     = New(codes.InvalidParam, "invalid parameter")
	ErrUnauthenticated  = New(codes.Unauthenticated, "authenticate failed")
	ErrPermissionDenied = New(codes.PermissionDenied, "permission denied")
	ErrExists           = New(codes.Exists, "already exists")
	ErrNotExists        = New(codes.NotExists, "not exists")
	ErrVersionOld       = New(codes.VersionOld, "version too old")
)

type ServiceError struct {
	code codes.Code
	desc string
}

func New(code codes.Code, desc string) error {
	return &ServiceError{code: code, desc: desc}
}

func (e *ServiceError) Code() codes.Code {
	return e.code
}

func (e *ServiceError) Error() string {
	return e.desc
}
