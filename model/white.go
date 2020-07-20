package model

import (
	"api/pdo"
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"runtime"
)

type White struct {
	Key   string
	Value Values
}

type Values struct {
	Ips      []string
	Accounts []string
}

func (w *White) Get() {

	mongo := pdo.Mongo.Database("sycenter").Collection("server_whites")

	filter := bson.D{{"key", "white"}}
	err := mongo.FindOne(context.TODO(), filter).Decode(&w)
	if err != nil {
		fmt.Printf("err:%#v\n", err)
		return
	}
}

func (w *White) IsIpExist(ip string, isip chan bool) {

	for _, v := range w.Value.Ips {
		if v == ip {
			isip <- true
			runtime.Goexit()
		}
	}

	isip <- false

}

func (w *White) IsAccountExist(account string, isaccount chan bool) {

	for _, v := range w.Value.Accounts {
		if v == account {
			isaccount <- true
			runtime.Goexit()
		}
	}

	isaccount <- false

}

func (w *White) IsExist(user *User) bool {

	var isip = make(chan bool)
	var isaccount = make(chan bool)

	go w.IsIpExist(user.Ip, isip)

	go w.IsAccountExist(user.Account, isaccount)

	if <-isip || <-isaccount {
		return true
	}

	return false

}
