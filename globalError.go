package day4

import "net/http"

// GlobalError 全局异常结构体
type GlobalError struct {
	Status  int    `json:"-"`
	Code    int    `json:"code"`
	Message string `json:"message"`
}

//获取err的内容
func (err GlobalError) Error() string {
	return err.Message
}

const (
	NoReg           = 1001 //账号不存在请重新输入或注册
	InternalService = 1002 //内部服务错误
)

func NoRegError(s string) GlobalError {
	return GlobalError{
		Status:  http.StatusOK,
		Code:    NoReg,
		Message: s,
	}
}

//InternalServiceError   内部服务错误
func InternalServiceError(s string) GlobalError {
	return GlobalError{
		Status:  http.StatusOK,
		Code:    InternalService,
		Message: s,
	}
}
