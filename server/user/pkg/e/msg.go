package e

const (
	VALID                                = "valid"
	OK                                   = "ok"
	TYPE_ERROR                           = "type error"
	UPDATE_TOKEN_FAIL                    = "update token fail and token out of date"
	TOKEN_OUT_OF_DATE                    = "token out of date"
	USERNAME_NOT_EXIST_OR_PASSWORD_WRONG = "username or password wrong"
	PHONE_NOT_EXIST                      = "phone not exist"
	INFO_ERROR                           = "info not exist or wrong"
	BIND_JSON_ERROR                      = "bind json error"
	TOKEN_ERROR                          = "token error"
	MYSQL_ERROR                          = "mysql error"
	MYSQL_AFFAIR_ERROR                   = "mysql affair error"
	FILE_FORMAT_OR_SIZE_ERROR            = "file format or size error"
	UPLOAD_CHECK_FILE_ERROR              = "check file error"
	UPLOAD_SAVE_FILE_ERROR               = "save file error"
	ROUTING_NOT_EXIST                    = "router not exist"
	APP_ERROR                            = "platform error"
	OPERATING_ERROR                      = "operating error"
	NOT_LOGIN                            = "websocket user not login"
	SERVER_ERROR                         = "server error"
	COMMUNITY_ID_NOT_EXIST               = "communityID not exist"
	PERMISSION_DENY                      = "no permission"
	REPEAT_SUBMIT                        = "repeat submit"
	INFO_NOT_EXIST                       = "info not exist or id error"
)

var msgFlags = map[int]string{
	Valid:                           VALID,
	Ok:                              OK,
	TypeError:                       TYPE_ERROR,
	UpdateTokenFail:                 UPDATE_TOKEN_FAIL,
	TokenOutOfDate:                  TOKEN_OUT_OF_DATE,
	UserNameNotExistOrPasswordWrong: USERNAME_NOT_EXIST_OR_PASSWORD_WRONG,
	PhoneNotExist:                   PHONE_NOT_EXIST,
	InfoError:                       INFO_ERROR,
	TokenError:                      TOKEN_ERROR,
	MysqlError:                      MYSQL_ERROR,
	FileFormatOrSizeError:           FILE_FORMAT_OR_SIZE_ERROR,
	UploadCheckFileError:            UPLOAD_CHECK_FILE_ERROR,
	UploadSaveFileError:             UPLOAD_SAVE_FILE_ERROR,
	RoutingNotExist:                 ROUTING_NOT_EXIST,
	AppError:                        APP_ERROR,
	OperatingError:                  OPERATING_ERROR,
	NotLogin:                        NOT_LOGIN,
	ServerError:                     SERVER_ERROR,
	CommunityIDNotExist:             COMMUNITY_ID_NOT_EXIST,
	PermissionDeny:                  PERMISSION_DENY,
	RepeatSubmit:                    REPEAT_SUBMIT,
	InfoNotExist:                    INFO_NOT_EXIST,
}

// 获取错误代码对应的中文意思
func GetMsg(code int) string {
	msg, ok := msgFlags[code]
	if ok {
		return msg
	}

	return msgFlags[ERROR]
}
