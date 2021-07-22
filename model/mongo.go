package model

import "gopkg.in/mgo.v2"

type Message struct {
	UID        string `json:"uid"`        //用户id
	GoldNum    string `json:"goldnum"`    //金币数
	DiamondNum string `json:"diamondnum"` //钻石数
}

var Client *mgo.Collection

func GetClient() *mgo.Collection {
	mongo, _ := mgo.Dial("127.0.0.1")
	client := mongo.DB("mydb_gift").C("g_user")
	return client
}
