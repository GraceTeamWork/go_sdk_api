package model

import (
	"api/pdo"
	"context"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"log"
	"runtime"
	"strconv"
	"sync"
	"time"
)

type ServerData struct {
	Channel string    `bson:"channel"`
	Data    []*Server `bson:"data" json:"data"`
}

type Server struct {
	Default      int    `bson:"default" json:"default"`
	IsShow       int    `bson:"is_show" json:"is_show"`
	GroupId      int    `bson:"group_id" json:"group_id"`
	IsCreate     int    `bson:"is_create" json:"is_create,omitempty"`
	OpenTime     string `bson:"open_time" json:"open_time"`
	CombineTime  string `bson:"combine_time" json:"combine_time,omitempty"`
	ServerId     string `bson:"server_id" son:"server_id,omitempty"`
	Name         string `bson:"name" json:"name"`
	Flag         string `bson:"flag" json:"flag"`
	Sort         int    `bson:"sort" json:"sort"`
	Status       int    `bson:"status" json:"status"`
	Zone         int    `bson:"zone,,int" json:"zone,int"`
	Md5Key       string `bson:"md_5_key" json:"md_5_key,omitempty"`
	Domain       string `bson:"domain" json:"domain,omitempty"`
	ApiUrl       string `bson:"api_url" json:"api_url,omitempty"`
	MaxOnline    int    `bson:"max_online" json:"max_online"`
	ServerAddr   string `bson:"server_addr" json:"server_addr"`
	ServerPort   int    `bson:"server_port" son:"server_port"`
	ChatAddr     string `bson:"chat_addr" json:"chat_addr"`
	ChatPort     int    `bson:"chat_port" son:"chat_port"`
	ChannelFlag  string `bson:"channel_flag" json:"channel_flag,omitempty"`
	PackageFlags string `bson:"package_flags" json:"package_flags,omitempty"`
	CreatedTime  string `bson:"created_time" json:"created_time,omitempty"`
	Text         string `json:"test"`
	Sign         string `json:"sign"`
	Tick         int64  `json:"tick"`
}

type SingleServer struct {
	// Flag    string `gorm:"flag"`
	ApiUrl string `gorm:"api_url"`
}

var lock = sync.RWMutex{}

func (s Server) GetServerList(user *User) ([]*Server, error) {

	mongo := pdo.Mongo.Database("sycenter").Collection("server_list")

	var mogoData ServerData

	filter := bson.D{{"channel", user.Channel}}
	err := mongo.FindOne(context.TODO(), filter).Decode(&mogoData)
	if err != nil {
		//fmt.Printf("err:%#v\n", err)
		return nil, err
	}
	return mogoData.Data, nil
}

func DefaultFilter(user User, server []*Server) (data []*Server) {
	// 当前时间
	tick := time.Now().Unix()
	// 白名单
	var w White
	w.Get()
	whiteIsExist := w.IsExist(&user)
	for _, v := range server {
		// 校验包
		var packages []string
		err := json.Unmarshal([]byte(v.PackageFlags), &packages)
		if err != nil {
			log.Fatal(err)
			continue
		}
		var packChan = make(chan bool)
		go checkPackage(packages, strconv.Itoa(user.Package), packChan)
		// 过滤未开服且不显示的服
		optime, err := time.Parse("2006-01-02 15:04:05", v.OpenTime)
		if err != nil {
			fmt.Printf("解析时间错误:%#v\n", err)
		}
		//白名单
		if v.Status < 0 {
			if whiteIsExist {
				v.Status = 1
			}
		}
		// 包
		if !<-packChan {
			continue
		}
		if optime.Unix() > tick && v.IsShow == 0 {
			continue
		}

		if v.Status == -1 {
			v.Text = "维护中"
		} else {
			v.Text = v.OpenTime + " 火爆开启"
		}
		v.Sign = fmt.Sprintf("%x", md5.Sum([]byte(user.Account+strconv.Itoa(int(tick))+v.Md5Key)))
		v.Tick = tick
		// 记录大区
		zoneServer[strconv.Itoa(v.Zone)] = 1
		//SetZone(strconv.Itoa(v.Zone))

		// 置空，json时忽略
		v.IsCreate = 0
		v.CombineTime = ""
		v.ServerId = ""
		v.Md5Key = ""
		v.Domain = ""
		v.ApiUrl = ""
		v.ChannelFlag = ""
		v.PackageFlags = ""
		v.CreatedTime = ""
		data = append(data, v)
	}
	return data
}

func checkPackage(packages []string, pack string, packChan chan bool) {
	for _, v := range packages {
		if v == pack {
			packChan <- true
			runtime.Goexit()
		}
	}
	packChan <- false
	return
}

//map 高并发读写必须加锁
func SetZone(zone string) {
	lock.Lock()
	defer lock.Unlock()
	zoneServer[zone] = 1
}

// 获取单服url_api

func (ss *SingleServer) getApiUrl(flag string) error {

	if err := pdo.Db.Select("sc.api_url").Table("xxwd_server as sb").Joins("LEFT join xxwd_server_config as sc on sc.flag = sb.cflag").Where("sb.flag = ?", flag).Find(&ss).Error; err != nil {
		return err
	}

	return nil
}
