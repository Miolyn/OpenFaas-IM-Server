package service

import (
	"OpenFaas-Connect/pkg/dto"
	"fmt"

	"github.com/imroc/req"
)

var (
	//proxyServerBind = "http://127.0.0.1:8083/check_token"
	proxyServerBind = "http://119.3.229.43:31112/function/proxy/check_token"
)

type RegisterService struct{}

func NewRegisterService() *RegisterService {
	return &RegisterService{}
}

//在logic服务器注册用户 代表上线
func (s *RegisterService) Register(token string) (form *dto.BindResponseJson, err error) {
	url := proxyServerBind
	authHeader := req.Header{
		"Content-Type": "application/json",
		"token":        token,
	}
	type BindReq struct {
		UserId   string `json:"user_id"`
		ConnId   int64  `json:"conn_id"`
		DeviceID int64  `json:"device_Id"`
	}
	resp, err := req.Get(url, authHeader, nil)
	if err != nil {
		return nil, fmt.Errorf("请求logic 服务器失败:%v", err)
	}
	res := new(dto.BindResponseJson)
	err = resp.ToJSON(res)
	if err != nil {
		return nil, fmt.Errorf("resp parse error :%v", err)
	}
	if res.Code != 0 {
		return nil, fmt.Errorf("logic server error:%v", res.Msg)
	}
	// TODO
	return res, nil
}

//在logic服务器注销用户 代表下线
func (s *RegisterService) UnRegister() (int64, error) {
	return 0, nil
}
