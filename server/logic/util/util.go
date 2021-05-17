package util

import (
	"OpenFaaS-Logic/st"
	"github.com/go-playground/validator/v10"
)

func CheckValidateForm(form interface{}) bool {
	validate := validator.New()
	err := validate.Struct(form)
	if err != nil {
		st.DebugWithFuncName(err)
		return true
	}
	return false
}
