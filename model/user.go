package model

import (
	"api/pdo"
	"context"
	"go.mongodb.org/mongo-driver/bson"
)

type User struct {
	CommonLogin
	NickName string
	Account  string `json:"openId"`
	Ip       string `json:"-"`
}

type UserInfo struct {
	Account  string `bson:"-"`
	Server   string `bson:"server"`
	Career   string `bson:"-"`
	Channel  string `bson:"-"`
	Level    string `bson:"level"`
	Package  string `bson:"-"`
	RoleId   string `bson:"-"`
	RoleName string `bson:"role_name"`
}

type Register struct {
	Account    string `gorm:"account"`
	Channel    string `gorm:"channel"`
	Package    string `gorm:"package"`
	System     string `gorm:"system"`
	Imei       string `gorm:"imei"`
	Idfa       string `gorm:"idfa"`
	CreateTime int64  `gorm:"create_time"`
	UpdateTime int64  `gorm:"create_time"`
}

func (u *User) GetUserInfo() map[string]interface{} {
	var userInfo []UserInfo
	mongo := pdo.Mongo.Database("sycenter").Collection("level_info")
	filter := bson.D{{"account", u.Account}, {"channel", u.Channel}, {"package", u.Package}}
	err := mongo.FindOne(context.TODO(), filter).Decode(&userInfo)
	if err != nil {
		return nil
	}
	var returnData = make(map[string]interface{})
	for _, v := range userInfo {
		var tmp = make(map[string]interface{})
		tmp[v.RoleName] = v.RoleName
		tmp[v.Level] = v.Level
		returnData[v.Server] = tmp
	}
	return returnData
}
