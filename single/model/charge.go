package model

import (
	"api/pdo"
	"context"
	"crypto/md5"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	Key = "kevvetr5z2t0jzgnwse3w64akwgrzl3mxod"
)

type Order struct {
	Server       string  `bson:"server" form:"server"`
	Channel      string  `bson:"channel" form:"channel"`
	Package      int     `bson:"package" form:"package"`
	Id           string  `bson:"id" form:"id"`
	Money        float64 `bson:"money" form:"money"`
	MoneyType    int     `bson:"money_type" form:"money_type"`
	Gold         int     `bson:"gold" form:"gold"`
	ConfId       int     `bson:"conf_id" form:"conf_id"`
	Account      string  `bson:"account" form:"account"`
	RoleId       int     `bson:"role_id" form:"role_id"`
	RoleLevel    int     `bson:"role_level" form:"role_level"`
	RechargeTick float64 `bson:"recharge_tick" form:"recharge_tick"`
	Tick         string  `bson:"-" form:"tick"`
	Sign         string  `bson:"-" form:"sign"`
}

func (ot *Order) AddOrder() error {
	if ot.Id == "" {
		return errors.New("订单有误")
	}

	if ot.Server == "" {
		return errors.New("订单有误")
	}

	if ot.Channel == "" {
		return errors.New("订单有误")
	}

	if ot.Package == 0 {
		return errors.New("订单有误")
	}

	if ot.Account == "" {
		return errors.New("订单有误")
	}

	if ot.MoneyType <= 0 || ot.MoneyType >= 8 {
		return errors.New("金钱类型错误")
	}
	if ot.Money <= 0 {
		return errors.New("金额错误")
	}
	if ot.ConfId <= 0 {
		return errors.New("充值类型错误")
	}

	if ot.Tick == "" {
		return errors.New("订单有误")
	}

	if ot.Sign == "" {
		return errors.New("订单有误")
	}

	sign := fmt.Sprintf("%x", md5.Sum([]byte(ot.Tick+Key)))

	if sign != ot.Sign {
		return errors.New("签名错误")
	}

	collection := pdo.Mongo.Database("sycenter").Collection("order")

	filter := bson.M{"id": ot.Id}                  //条件
	update := bson.M{"$set": ot}                   //更新数据
	var UpInert = options.Update().SetUpsert(true) //没有集合时插入集合
	Result, err := collection.UpdateOne(context.TODO(), filter, update, UpInert)
	if err != nil {
		return err
	}
	fmt.Printf("res:%v\n", Result)

	return nil

}
