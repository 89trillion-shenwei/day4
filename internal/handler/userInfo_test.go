package handler

import (
	model2 "day4/internal/model"
	"day4/internal/service"
	"day4/internal/util"
	"encoding/json"
	"fmt"
	"github.com/golang/protobuf/proto"
	"testing"
)

func TestLogin(t *testing.T) {
	client := model2.GetClient()
	str, _ := Login("1626945193510426000", client)
	var data []byte = []byte(str)
	message := model2.Message{}
	err := json.Unmarshal(data, &message)
	if err != nil {
		t.Log(err.Error())
		return
	}
	if message.UID == "1626945193510426000" && message.GoldNum == "0" && message.DiamondNum == "0" {
		t.Log("success")
		return
	} else {
		t.Log("failed")
		t.Error("failed")
	}
}

func TestRegister(t *testing.T) {
	client := model2.GetClient()
	uid, _ := Register(client)
	str, _ := Login(uid, client)
	var data = []byte(str)
	message := model2.Message{}
	err := json.Unmarshal(data, &message)
	if err != nil {
		t.Log(err.Error())
		return
	}
	if message.UID == uid && message.GoldNum == "0" && message.DiamondNum == "0" {
		t.Log("success")
		return
	} else {
		t.Log("failed")
		t.Error("failed")
	}
}

func TestUpdate(t *testing.T) {
	util.Init()
	//data :=make([]byte,0)
	client := model2.GetClient()
	uid, err := Register(client)
	if err != nil {
		fmt.Println(err)
	}
	user := service.User{}
	user.UserName = uid
	data, err1 := Update(user, "5d22480d")
	if err1 != nil {
		t.Log(err1)
	}
	ge := model2.GeneralReward{}
	proto.Unmarshal(data, &ge)
	if err1 != nil {
		fmt.Println(err1)
	}
	if ge.Changes[1001] == 4 && ge.Changes[1002] == 8 {
		t.Log("success")
		return
	} else {
		t.Log("failed")
		t.Error("failed")
	}
}
