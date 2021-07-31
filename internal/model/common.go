package model

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/garyburd/redigo/redis"
	"gopkg.in/mgo.v2/bson"
)

var Data Mess

// List 物品列表
type List struct {
	Name   string //物品名
	Amount string // 物品数量
}

// GetList 领取信息列表
type GetList struct {
	GetorName string //领取人用户名
	GetTime   string //领取时间
}

// Message redis存储的礼品码信息
type Message struct {
	Description string //礼品描述
	CodeType    string //礼品码类型
	List        []List //礼品内容列表（物品，数量）
	ValidPeriod string //有效期
	GiftCode    string //礼品码
	CanGetUser  string //允许领取用户
	Creator     string //创建者账号
	CreatTime   string //创建时间
}

// Mess redis存储的领取信息
type Mess struct {
	AvailableTimes string    //可领取次数
	ReceivedTimes  string    //已领取次数
	GiftCode       string    //礼品码
	key            string    //计数
	GetList        []GetList //领取列表
}

//结构体转json
func Struct2json(me interface{}) []byte {
	byts, _ := json.Marshal(me)
	return byts
}

//json转结构体
func Json2struct(byts []byte, message *Message) bool {
	err := json.Unmarshal(byts, message)
	if err != nil {
		return false
	}
	return true

}

//json转结构体
func Json2struct1(byts []byte, message *Mess) bool {
	err := json.Unmarshal(byts, message)
	if err != nil {
		return false
	}
	return true

}

// CheckKey 判断礼品码是否存在
func CheckKey(key string) bool {
	c := RedisPool.Get()
	defer c.Close()
	exist, err := redis.Bool(c.Do("EXISTS", key))
	if err != nil {
		return false
	} else {
		return exist
	}
}

// CheckId 判断用户是否存在
func CheckId(uid string) bool {
	c := Client
	message := Message{}
	err := c.Find(bson.M{"uid": uid}).One(&message)
	if err != nil {
		fmt.Println("未找到此uid，请注册")
		return false
	}
	return true
}

//判断用户是否已经领取过礼品
func FindUser(list []GetList, name string) bool {
	for _, item := range list {
		if item.GetorName == name {
			return true
		}
	}
	return false
}

// RandomCode 根据时间戳生成uid
func RandomCode() string {
	code := time.Now().UnixNano()
	str := strconv.FormatInt(code, 10)
	return str
}

// String2Int32 字符串转int32
func String2Int32(str string) int32 {
	in, _ := strconv.ParseInt(str, 10, 32)
	return int32(in)
}

// String2UInt32 字符串转uint32
func String2UInt32(str string) uint32 {
	in, _ := strconv.ParseInt(str, 10, 32)
	return uint32(in)
}

// String2UInt64 字符串转uint64
func String2UInt64(str string) uint64 {
	in, _ := strconv.Atoi(str)
	return uint64(in)
}
