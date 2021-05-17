package e

const (
	ERROR = 500
)

const baseIndex = 1

const (
	Valid = baseIndex + iota
	Ok
	TypeError
	UpdateTokenFail

	TokenOutOfDate
	TokenError
	MysqlError

	UserNameNotExistOrPasswordWrong
	PhoneNotExist
	InfoError

	FileFormatOrSizeError

	UploadCheckFileError
	UploadSaveFileError

	CommunityIDNotExist
	PermissionDeny
	RoutingNotExist

	InfoNotExist

	RepeatSubmit
	AppError
	OperatingError
	NotLogin
	ServerError
)
