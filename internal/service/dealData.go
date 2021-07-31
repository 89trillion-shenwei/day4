package service

import (
	"day4/internal"
	"day4/internal/model"
	"fmt"
	"github.com/garyburd/redigo/redis"
	"log"
	"strconv"
	"time"
)

// Creator 管理员信息
type Creator struct {
	CreaName string //管理员账号
}

// User 用户信息
type User struct {
	UserName string //用户账号
}

// StrUpdate 用户领取礼品时更新数据库，增加领取人列表，修改可领取次数和已领取次数，返回礼品列表
func (User *User) StrUpdate(key string) ([]model.List, error) {
	c1 := model.RedisPool.Get()
	c2 := model.RedisPool1.Get()
	defer c1.Close()
	defer c2.Close()
	//查询数据
	res1, err1 := redis.String(c1.Do("Get", key))
	res2, err2 := redis.String(c2.Do("Get", key))
	if err1 != nil || err2 != nil {
		return nil, internal.InternalServiceError(err1.Error() + err2.Error())
	} else {
		var byts1, byts2 []byte
		byts1 = []byte(res1) //礼品码信息
		byts2 = []byte(res2) //领取信息
		mess := model.Mess{}
		message := model.Message{}
		//json转结构体成功
		if model.Json2struct(byts1, &message) && model.Json2struct1(byts2, &mess) {
			fmt.Println("success")
		}
		model.Data = mess
		//如果礼品可领取次数为0，退出
		if mess.AvailableTimes == "0" {
			return nil, internal.NoGiftError("礼品已领完")
		}
		//指定用户一次性消耗
		if message.CodeType == "1" {
			//如果用户是指定用户
			if message.CanGetUser == User.UserName {
				getList := new(model.GetList)
				//用户名
				getList.GetorName = User.UserName
				//领取时间
				getList.GetTime = time.Now().Format("2006-01-02 15:04:05")
				mess.GetList = append(mess.GetList, *getList)
				//可领取次数
				av, _ := strconv.Atoi(mess.AvailableTimes)
				av -= 1
				//已领取次数
				re, _ := strconv.Atoi(mess.ReceivedTimes)
				re += 1
				mess.ReceivedTimes = strconv.Itoa(re)
				mess.AvailableTimes = strconv.Itoa(av)
				//提交更改后的数据
				_, err := c2.Do("SET", key, model.Struct2json(mess))
				if err != nil {
					return nil, internal.InternalServiceError(err.Error())
				} else {
					fmt.Println("set ok.")
				}
				//incr key
				_, err1 := redis.Int64(c2.Do("INCR", "key"))
				if err1 != nil {
					log.Println("INCR failed:", err)
					return nil, internal.InternalServiceError(err1.Error())
				}
			} else {
				return nil, internal.NoCanGetUserError("非指定用户")
			}
		} else if message.CodeType == "2" { //不指定用户,限制次数
			getList := new(model.GetList)
			//用户名
			getList.GetorName = User.UserName
			//领取时间
			getList.GetTime = time.Now().Format("2006-01-02 15:04:05")
			//判断该用户是否已经使用过该礼品码
			if model.FindUser(mess.GetList, User.UserName) {
				return nil, internal.UserHasEeceivedError("你已使用过该礼品码")
			}
			mess.GetList = append(mess.GetList, *getList)
			//可领取次数
			av, _ := strconv.Atoi(mess.AvailableTimes)
			av -= 1
			//已领取次数
			re, _ := strconv.Atoi(mess.ReceivedTimes)
			re += 1
			mess.ReceivedTimes = strconv.Itoa(re)
			mess.AvailableTimes = strconv.Itoa(av)
			//提交更改后的数据
			_, err := c2.Do("SET", key, model.Struct2json(mess))
			if err != nil {
				return nil, internal.InternalServiceError(err.Error())
			} else {
				fmt.Println("set ok.")
			}
			//incr key
			_, err1 := redis.Int64(c2.Do("INCR", "key"))
			if err1 != nil {
				log.Println("INCR failed:", err)
				return nil, internal.InternalServiceError(err1.Error())
			}
		} else { //不限制用户，不限制次数
			getList := new(model.GetList)
			//用户名
			getList.GetorName = User.UserName
			//领取时间
			getList.GetTime = time.Now().Format("2006-01-02 15:04:05")
			//判断该用户是否已经使用过该礼品码
			if model.FindUser(mess.GetList, User.UserName) {
				return nil, internal.UserHasEeceivedError("你已使用过该礼品码")
			}
			mess.GetList = append(mess.GetList, *getList)
			//已领取次数
			re, _ := strconv.Atoi(mess.ReceivedTimes)
			re += 1
			mess.ReceivedTimes = strconv.Itoa(re)
			mess.AvailableTimes = "999999"
			//提交更改后的数据
			_, err := c2.Do("SET", key, model.Struct2json(mess))
			if err != nil {
				return nil, internal.InternalServiceError(err.Error())
			} else {
				fmt.Println("set ok.")
			}
			//incr key
			_, err1 := redis.Int64(c2.Do("INCR", "key"))
			if err1 != nil {
				log.Println("INCR failed:", err)
				return nil, internal.InternalServiceError(err1.Error())
			}
		}
		//返回礼品内容
		return message.List, nil
	}
}

// ReturnBack 回退函数
func ReturnBack(key string) error {
	c := model.RedisPool1.Get()
	_, err := c.Do("SET", key, model.Struct2json(model.Data))
	if err != nil {
		return internal.InternalServiceError("回调失败")
	} else {
		fmt.Println("已回调")
		return nil
	}
}
