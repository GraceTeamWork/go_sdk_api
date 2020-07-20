package model

import (
	"api/common"
	"api/pdo"
	"crypto/md5"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"sync"
	"time"
)

/*
 *	When I wrote this, only God and I understood what I was doing
 *	Now, God only knows
 *  Author: liuhuacong
 *  Time:	2020-07-10
 */
const (
	Key = "kevvetr5z2t0jzgnwse3w64akwgrzl3mxod"
)

type Orde interface {
	CheakParam() error
	CheakSign() error
	Order() error
}

var ow sync.WaitGroup

type Order struct {
	OrderNum   string  `gorm:"order_num"`
	CpOrderNum string  `gorm:"cp_order_num"`
	Channel    string  `gorm:"channel"`
	Package    int     `gorm:"Package"`
	Server     string  `gorm:"server"`
	Money      float64 `gorm:"money"`
	MoneyType  int     `gorm:"money_type"`
	Coupon     float64 `gorm:"coupon"`
	Gold       int     `gorm:"gold"`
	PayWay     string  `gorm:"pay_way"`
	ConfId     int     `gorm:"conf_id"`
	Account    string  `gorm:"account"`
	RoleId     int     `gorm:"role_id"`
	Commodity  string  `gorm:"commodity"`
	RoleName   string  `gorm:"role_name"`
	RoleLevel  int     `gorm:"role_level"`
	Sync       int     `gorm:"sync"`
	SyncTime   int64   `gorm:"sync_time"`
	SyncMsg    string  `gorm:"sync_msg"`
	Currency   string  `gorm:"currency"`
	System     string  `gorm:"system"`
	First      int     `gorm:"first"`
	RoleFirst  int     `gorm:"role_first"`
	IsTest     int     `gorm:"is_test"`
	Ip         string  `gorm:"ip"`
	Mac        string  `gorm:"mac"`
	Imei       string  `gorm:"imei"`
	Idfa       string  `gorm:"idfa"`
	Extra      string  `gorm:"extra"`
	CreateTime int64   `gorm:"create_time"`
}

type HttpRespone struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

func AddOrder(o Orde) error {

	if err := o.CheakSign(); err != nil {
		return err
	}

	if err := o.CheakParam(); err != nil {
		return err
	}

	if err := o.Order(); err != nil {
		return err
	}

	return nil
}

func (ot *Order) AddOrder() error {

	// 判断订单是否存在
	// 测试
	midOder := CountNum()()
	ot.OrderNum = ot.OrderNum + midOder

	fmt.Println(ot.OrderNum)
	if err := ot.GetOrderInfo(); err == nil {
		return errors.New("订单已存在")
	}
	if ot.Money <= 0 || ot.Money > 10000000 {
		return errors.New("订单金额有误")
	}
	if ot.Gold < 0 {
		return errors.New("元宝有误")
	}
	if ot.Gold == 0 {
		ot.Gold = int(ot.Money * 10)
	}

	if ot.Server == "" {
		return errors.New("服务器有误")
	}
	if ot.Channel == "" {
		return errors.New("渠道有误")
	}
	if ot.Package == 0 {
		return errors.New("包有误")
	}
	var ch = make(chan struct{}, 3)
	ow.Add(3)
	go ot.IsFirst(ch)
	go ot.IsRoleFirst(ch)
	go ot.system(ch)
	ow.Wait()
	close(ch)
	for {
		_, isOk := <-ch
		if !isOk {
			break
		}
	}

	if err := ot.Insert(); err != nil {
		return err
	}

	ot.Syncs()
	return nil
}

func (ot *Order) GetOrderInfo() error {

	if err := pdo.Db.Model(&ot).Where("order_num = ?", ot.OrderNum).First(&ot).Error; err != nil {
		return err
	}
	return nil
}

// 账号首充
func (ot *Order) IsFirst(ch chan struct{}) {
	defer ow.Done()
	var count int
	if err := pdo.Db.Model(&ot).Where("account = ?", ot.Account).Count(&count).Error; err != nil {
		// ot.First = 1
		ch <- struct{}{}
		return
	}
	ot.First = count + 1
	ch <- struct{}{}
	return
}

// 角色首充
func (ot *Order) IsRoleFirst(ch chan struct{}) {
	defer ow.Done()
	var order Order
	if err := pdo.Db.Model(&ot).Where("role_id = ?", ot.RoleId).First(&order).Error; err != nil {
		ot.RoleFirst = 1
		ch <- struct{}{}
		return
	}
	ot.RoleFirst = 0
	ch <- struct{}{}
	return
}

// 充值系统
func (ot *Order) system(ch chan struct{}) {
	defer ow.Done()
	var reg Register
	if err := pdo.Db.Model(&reg).Where("account = ?", ot.Account).First(&reg).Error; err != nil {
		ot.System = "none"
		ch <- struct{}{}
		return
	}
	ot.System = reg.System
	ch <- struct{}{}
	return
}

// 插入订单
func (ot *Order) Insert() error {
	if err := pdo.Db.Model(&ot).Create(&ot).Error; err != nil {
		return err
	}
	return nil
}

// 更新订单
func (ot *Order) Update() error {
	if err := pdo.Db.Model(&ot).Where("order_num = ?", ot.OrderNum).Update(&ot).Error; err != nil {
		return err
	}
	return nil
}

// 同步订单
func (ot *Order) Syncs() {
	// 获取单服url_api
	var server SingleServer
	err := server.getApiUrl("cs001")
	jsonParam, _ := json.Marshal(ot)
	if err != nil {
		ot.Sync = 2
		ot.SyncMsg = "服务器不正确"
		ot.SyncTime = time.Now().Unix()
		ot.Update()
		common.SyncLog(ot.Channel, string(jsonParam)+" error:"+"服务器不正确")
		return

	} else {
		server.ApiUrl = "http://127.0.0.1:9003"

		var param = url.Values{}

		param.Set("server", ot.Server)
		param.Set("channel", ot.Channel)
		param.Set("package", strconv.Itoa(ot.Package))
		param.Set("id", ot.OrderNum)

		param.Set("money", strconv.FormatFloat(ot.Money, 'E', 2, 64))
		param.Set("money_type", strconv.Itoa(ot.MoneyType))
		param.Set("gold", strconv.Itoa(ot.Gold))
		param.Set("conf_id", strconv.Itoa(ot.ConfId))

		param.Set("role_id", strconv.Itoa(ot.RoleId))
		param.Set("account", ot.Account)
		param.Set("role_level", strconv.Itoa(ot.RoleLevel))
		param.Set("recharge_tick", strconv.Itoa(int(ot.CreateTime)))

		tick := strconv.Itoa(int(time.Now().Unix()))
		param.Set("tick", tick)
		param.Set("sign", fmt.Sprintf("%x", md5.Sum([]byte(tick+Key))))
		var ch = make(chan struct{})
		resp, err := http.PostForm(server.ApiUrl+"/charge", param)
		if err != nil {
			go func() {
				common.SyncLog(ot.Channel, string(jsonParam)+" error:"+err.Error())
				ch <- struct{}{}
			}()
			ot.Sync = 2
			ot.SyncTime = time.Now().Unix()
			ot.SyncMsg = err.Error()
			ot.Update()
			log.Printf("err:%v", err)
			<-ch
			return
		}
		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			go func() {
				common.SyncLog(ot.Channel, string(jsonParam)+" error:"+err.Error())
				ch <- struct{}{}
			}()
			ot.Sync = 2
			ot.SyncMsg = err.Error()
			ot.SyncTime = time.Now().Unix()
			ot.Update()
			log.Printf("err:%v", err)
			<-ch
			return
		}
		var res HttpRespone
		err = json.Unmarshal(body, &res)
		if err != nil {
			log.Printf("err:%v", err)
		}

		go func() {
			httpRespone, _ := json.Marshal(body)
			common.SyncLog(ot.Channel, string(jsonParam)+" HttpRespone:"+string(httpRespone))
			ch <- struct{}{}
		}()
		if res.Code == 200 {
			ot.Sync = 1
		} else {
			ot.Sync = 4
		}
		ot.SyncMsg = res.Msg
		ot.SyncTime = time.Now().Unix()
		ot.Update()
		<-ch
		return
	}

}
