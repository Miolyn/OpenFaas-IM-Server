package dto

type ResponseJson struct {
	Msg   string      `json:"msg"`
	Code  int         `json:"code"`
	Data  interface{} `json:"data"`
	Page  int64       `json:"page"`
	Size  int64       `json:"size"`
	Total int64       `json:"total"`
}

type FORM struct {
	Username string `json:"username"`
	UserId   string `json:"user_id"`
}

type BindResponseJson struct {
	Msg   string `json:"msg"`
	Code  int    `json:"code"`
	Data  FORM   `json:"data"`
	Page  int64  `json:"page"`
	Size  int64  `json:"size"`
	Total int64  `json:"total"`
}

// code
const (
	RESPONSE_OK_CODE   = 0
	RESPONSE_FAIL_CODE = -1
)
