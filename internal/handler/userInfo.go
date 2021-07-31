package handler

import (
	"day4/internal"
	message2 "day4/internal/message"
	model2 "day4/internal/model"
	"day4/internal/service"
	"encoding/json"
	"fmt"
	"github.com/golang/protobuf/proto"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"strconv"
	"sync"
)

var (
	wg       sync.WaitGroup
	lockChan = make(chan struct{}, 1)
)

// 如果lockChan中为空则阻塞
func getLock() {
	<-lockChan
}

// 重新填充lockChan
func releaseLock() {
	lockChan <- struct{}{}
}

//登录
func Login(uid string, client *mgo.Collection) (string, error) {
	message := model2.Message1{}
	err := client.Find(bson.M{"uid": uid}).One(&message)
	if err != nil {
		fmt.Println("未找到此uid，请注册")
		return "", internal.NoRegError("未找到此uid，请注册,登录失败,错误原因为：" + err.Error())
	}
	//登录成功，返回数据库数据
	message1 := model2.Message1{}
	err1 := client.Find(bson.M{"uid": uid}).One(&message1)
	if err1 != nil {
		return "", internal.InternalServiceError("根据uid查询数据失败，错误原因为：" + err.Error())
	}
	jsons, _ := json.Marshal(message1)
	return string(jsons), nil
}

//注册
func Register(client *mgo.Collection) (string, error) {
	message := model2.Message1{}
	message.UID = model2.RandomCode()
	message.DiamondNum = "0"
	message.GoldNum = "0"
	err := client.Insert(message)
	if err != nil {
		return "", internal.InternalServiceError("注册插入数据时出现错误，错误原因为：" + err.Error())
	}
	return message.UID, nil
}

// TestUpdate 更新储存uid的数据
func testUpdate(client *mgo.Collection, uid, goldnum, diamondnum string) error {
	//client :=model.SetConnect()
	selector := bson.M{"uid": uid}
	update := bson.M{"$set": bson.M{"goldnum": goldnum, "diamondnum": diamondnum}}
	err := client.Update(selector, update)
	if err != nil {
		fmt.Println(err)
		return err
	}
	fmt.Println("success")
	return nil
}

// Update 用户领取礼品，更新数据
func Update(user service.User, key string) ([]byte, error) {
	if !model2.CheckId(user.UserName) {
		return nil, internal.NoRegError("用户未注册")
	}
	if model2.CheckKey(key) {
		//上锁
		releaseLock()
		lists, err := user.StrUpdate(key)
		if err != nil {
			//解锁
			getLock()
			return nil, err
		}
		//声明接收数据的结构体
		general := message2.GeneralReward{}
		//变化量
		change := make(map[uint32]uint64)
		//余额
		balance := make(map[uint32]uint64)
		//终值
		counter := make(map[uint32]uint64)
		//储存领取物品id和数目
		for i := 0; i < len(lists); i++ {
			change[model2.String2UInt32(lists[i].Name)] = model2.String2UInt64(lists[i].Amount)
		}
		message := model2.Message1{}
		general.Code = model2.String2Int32(user.UserName) //用户uid
		general.Msg = "用户uid为" + user.UserName + ",金币id为1001，钻石id为1002"
		//查询当前mongo数据库中的用户金币和钻石数
		model2.Client.Find(bson.M{"uid": user.UserName}).One(&message)
		balance[1001] = model2.String2UInt64(message.GoldNum)
		balance[1002] = model2.String2UInt64(message.DiamondNum)
		//计算变化后的值
		counter[1001] = change[1001] + balance[1001]
		counter[1002] = change[1002] + balance[1002]
		fmt.Println(counter[1001])
		fmt.Println(counter[1002])
		//储存到mongo数据库中
		client := model2.GetClient()
		err1 := testUpdate(client, user.UserName, strconv.FormatUint(counter[1001], 10), strconv.FormatUint(counter[1002], 10))
		if err1 != nil {

			//回调
			if service.ReturnBack(key) != nil {
				//解锁
				getLock()
				return nil, service.ReturnBack(key)
			}
			//解锁
			getLock()
			return nil, internal.InternalServiceError("mongo操作失败，已回调")
		}
		general.Changes = change
		general.Balance = balance
		general.Counter = counter
		general.Ext = ""
		data, err2 := proto.Marshal(&general)
		if err2 != nil {
			//解锁
			getLock()
			return nil, internal.InternalServiceError("转码错误" + err2.Error())
		}
		//解锁
		getLock()
		return data, nil
	} else {
		return nil, internal.NoKeyError("礼品码不存在")
	}

}
