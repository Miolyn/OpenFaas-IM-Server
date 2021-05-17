package dto

type ResponseJson struct {
	Msg   string      `json:"msg"`
	Code  int         `json:"code"`
	Data  interface{} `json:"data"`
	Page  int64       `json:"page"`
	Size  int64       `json:"size"`
	Total int64       `json:"total"`
}

// code
const (
	RESPONSE_OK_CODE   = 0
	RESPONSE_FAIL_CODE = -1
)
