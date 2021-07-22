package ctrl

import (
	"day4/internal"
	"day4/internal/handler"
	"day4/internal/service"
	"day4/model"
	"github.com/gin-gonic/gin"
	"net/http"
)

//登录api
func LoginApi(c *gin.Context) (string, error) {
	//message :=model.Message{}
	uid := c.PostForm("uid")
	client := model.GetClient()
	json, err := handler.Login(uid, client)
	if err != nil {
		return "", err
	}
	return json, nil
}

//注册api
func RegisterApi(c *gin.Context) (string, error) {
	client := model.GetClient()
	uid, err := handler.Register(client)
	if err != nil {
		return "", err
	}
	return uid, nil
}

//领取礼品api
func ReceiveGiftsApi(c *gin.Context) ([]byte, error) {
	user := new(service.User)
	key := c.PostForm("key")
	if key == "" {
		return nil, internal.IsEmptyError("礼品码不能为空")
	}
	if len(key) != 8 {
		return nil, internal.LenFalseError("礼品码不合法")
	}
	user.UserName = c.PostForm("username")
	if user.UserName == "" {
		return nil, internal.IsEmptyError("用户名不能为空")
	}
	re, err := handler.Update(*user, key)
	if err != nil {
		return nil, internal.InternalServiceError(err.Error())
	}
	return re, nil
}

type Api1 func(c *gin.Context) (string, error)
type Api2 func(c *gin.Context) ([]byte, error)

func ReturnData(api Api1) gin.HandlerFunc {
	return func(c *gin.Context) {
		data, err := api(c)
		if err != nil {
			globalError := err.(internal.GlobalError)
			c.JSON(globalError.Status, globalError)
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"查询结果": data,
		})
	}
}

func ReturnProto(api Api2) gin.HandlerFunc {
	return func(c *gin.Context) {
		data, err := api(c)
		if err != nil {
			globalError := err.(internal.GlobalError)
			c.JSON(globalError.Status, data)
			return
		}
		c.JSON(http.StatusOK, data)
	}
}
