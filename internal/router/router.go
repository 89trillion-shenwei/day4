package router

import (
	"day4/internal/ctrl"
	"github.com/gin-gonic/gin"
)

func MongoRouter() *gin.Engine {
	r := gin.Default()
	//登录
	r.POST("/login", ctrl.ReturnData(ctrl.LoginApi))
	//注册
	r.POST("/register", ctrl.ReturnData(ctrl.RegisterApi))

	r.POST("/receiveGifts", ctrl.ReturnProto(ctrl.ReceiveGiftsApi))
	return r
}
