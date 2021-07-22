package handler

import (
	"day4"
	"day4/internal"
	"day4/internal/service"
	"day4/model"
	"encoding/json"
	"fmt"
	"github.com/golang/protobuf/proto"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"strconv"
	"time"
)

func Login(uid string, client *mgo.Collection) (string, error) {
	message := model.Message{}
	err := client.Find(bson.M{"uid": uid}).One(&message)
	if err != nil {
		fmt.Println("未找到此uid，请注册")
		return "", day4.NoRegError("未找到此uid，请注册,登录失败,错误原因为：" + err.Error())
	}
	//登录成功，返回数据库数据
	message1 := model.Message{}
	err1 := client.Find(bson.M{"uid": uid}).One(&message1)
	if err1 != nil {
		return "", internal.InternalServiceError("根据uid查询数据失败，错误原因为：" + err.Error())
	}
	jsons, _ := json.Marshal(message1)
	return string(jsons), nil
}

//注册
func Register(client *mgo.Collection) (string, error) {
	message := model.Message{}
	message.UID = RandomCode()
	message.DiamondNum = "0"
	message.GoldNum = "0"
	err := client.Insert(message)
	if err != nil {
		return "", internal.InternalServiceError("注册插入数据时出现错误，错误原因为：" + err.Error())
	}
	return message.UID, nil
}

// TestUpdate 更新储存uid的数据
func testUpdate(client *mgo.Collection, uid, goldnum, diamondnum string) {
	//client :=model.SetConnect()
	selector := bson.M{"uid": uid}
	update := bson.M{"$set": bson.M{"goldnum": goldnum, "diamondnum": diamondnum}}
	err := client.Update(selector, update)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("success")
}

// Update 用户领取礼品，更新数据
func Update(user service.User, key string) ([]byte, error) {
	if service.CheckKey(key) {
		lists, err := user.StrUpdate(key)
		if err != nil {
			return nil, internal.InternalServiceError("更新失败" + err.Error())
		}
		//声明接收数据的结构体
		general := model.GeneralReward{}
		//变化量
		change := make(map[uint32]uint64)
		//余额
		balance := make(map[uint32]uint64)
		//终值
		counter := make(map[uint32]uint64)
		//储存领取物品id和数目
		for i := 0; i < len(lists); i++ {
			change[String2UInt32(lists[i].Name)] = String2UInt64(lists[i].Amount)
		}
		message := model.Message{}
		general.Code = String2Int32(user.UserName) //用户uid
		general.Msg = "用户uid为" + user.UserName + ",金币id为1001，钻石id为1002"
		//查询当前mongo数据库中的用户金币和钻石数
		model.Client.Find(bson.M{"uid": user.UserName}).One(&message)
		balance[1001] = String2UInt64(message.GoldNum)
		balance[1002] = String2UInt64(message.DiamondNum)
		//计算变化后的值
		counter[1001] = change[1001] + balance[1001]
		counter[1002] = change[1002] + balance[1002]
		fmt.Println(counter[1001])
		fmt.Println(counter[1002])
		//储存到mongo数据库中
		client := model.GetClient()
		testUpdate(client, user.UserName, strconv.FormatUint(counter[1001], 10), strconv.FormatUint(counter[1002], 10))
		general.Changes = change
		general.Balance = balance
		general.Counter = counter
		general.Ext = ""
		data, err1 := proto.Marshal(&general)
		if err1 != nil {
			return nil, internal.InternalServiceError("转码错误" + err1.Error())
		}
		return data, nil
	} else {
		return nil, internal.NoKeyError("礼品码不存在")
	}

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

func String2UInt32(str string) uint32 {
	in, _ := strconv.ParseInt(str, 10, 32)
	return uint32(in)
}

func String2UInt64(str string) uint64 {
	in, _ := strconv.Atoi(str)
	return uint64(in)
}
