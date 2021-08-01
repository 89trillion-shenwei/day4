package internal

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
	NoKey           = 1001 //礼品码不存在
	UserHasEeceived = 1002 //不可重复领取
	NoGift          = 1003 //礼品全部领完
	IsEmpty         = 1004 //参数不能为空
	LenFalse        = 1005 //礼品码不合法
	InternalService = 1006 //内部服务错误
	NoReg           = 1007 //账号不存在请重新输入或注册
	NoCanGetUser    = 1008 //非指定用户
)

// NoKeyError  礼品码不存在
func NoKeyError(message string) GlobalError {
	return GlobalError{
		Status:  http.StatusOK,
		Code:    NoKey,
		Message: message,
	}
}

// UserHasEeceivedError 不可重复领取
func UserHasEeceivedError(message string) GlobalError {
	return GlobalError{
		Status:  http.StatusOK,
		Code:    UserHasEeceived,
		Message: message,
	}
}

// NoGiftError  礼品全部领完
func NoGiftError(message string) GlobalError {
	return GlobalError{
		Status:  http.StatusOK,
		Code:    NoGift,
		Message: message,
	}
}

//IsEmptyError  参数不能为空
func IsEmptyError(message string) GlobalError {
	return GlobalError{
		Status:  http.StatusOK,
		Code:    IsEmpty,
		Message: message,
	}
}

//LenFalseError  礼品码不合法
func LenFalseError(message string) GlobalError {
	return GlobalError{
		Status:  http.StatusOK,
		Code:    LenFalse,
		Message: message,
	}
}

//InternalServiceError   内部服务错误
func InternalServiceError(message string) GlobalError {
	return GlobalError{
		Status:  http.StatusOK,
		Code:    InternalService,
		Message: message,
	}
}

// NoRegError 未注册
func NoRegError(s string) GlobalError {
	return GlobalError{
		Status:  http.StatusOK,
		Code:    NoReg,
		Message: s,
	}
}

//NoCanGetUserError  非指定用户
func NoCanGetUserError(message string) GlobalError {
	return GlobalError{
		Status:  http.StatusOK,
		Code:    NoCanGetUser,
		Message: message,
	}
}
