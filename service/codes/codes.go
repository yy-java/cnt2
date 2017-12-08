package codes

type Code int32

const (
	OK               Code = 0
	ServerError      Code = 1
	InvalidParam     Code = 2
	Unauthenticated  Code = 3
	PermissionDenied Code = 4
	Exists           Code = 5
	NotExists        Code = 6
	VersionOld       Code = 7
)
