package model

import (
	"api/common"
	"api/pdo"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"sync"
	"time"
)

// var OL = sync.RWMutex{}
var mu = sync.Mutex{}
var num int

type OrderTmper interface {
	CheckParam() error
	Intert() error
}

type OrderTmp struct {
	GetOrderId int     `gorm:"-" form:"getOrderId"`
	OrderId    string  `gorm:"order_id" form:"order_id"`
	Server     string  `gorm:"server" form:"server"`
	Channel    string  `gorm:"channel" form:"channel"`
	Package    int     `gorm:"package" form:"package"`
	RoleId     int     `gorm:"role_id" form:"role_id"`
	Account    string  `gorm:"account" form:"account"`
	Money      float64 `gorm:"money" form:"money"`
	MoneyType  int     `gorm:"money_type" form:"money_type"`
	Status     int     `gorm:"status" form:"status"`
	ConfId     int     `gorm:"conf_id" form:"conf_id"`
	RoleLevel  int     `gorm:"role_level" form:"role_level"`
	Idfa       string  `gorm:"idfa" form:"idfa"`
	Imei       string  `gorm:"imei" form:"imei"`
	Mac        string  `gorm:"mac" form:"mac"`
	Ip         string  `gorm:"ip" form:"ip"`
	CreateTime int64   `gorm:"create_time"`
	UpdateTime int64   `gorm:"update_time"`
}

func GetOrderNum(ot OrderTmper) error {

	if err := ot.CheckParam(); err != nil {
		return err
	}

	if err := ot.Intert(); err != nil {
		return err
	}

	return nil
}

func (ot *OrderTmp) CheckParam() error {

	if ot.GetOrderId != 1 {
		return errors.New("参数错误")
	}
	if ot.Server == "" {
		return errors.New("服务器错误")
	}
	if ot.Channel == "" {
		return errors.New("渠道错误")
	}
	if ot.Package == 0 {
		return errors.New("包错误")
	}

	if ot.RoleId == 0 {
		return errors.New("角色id错误")
	}
	if ot.Account == "" {
		return errors.New("账号错误")
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
	if ot.RoleLevel <= 0 {
		return errors.New("角色等级错误")
	}

	return nil
}

func (ot *OrderTmp) Intert() error {

	ot.CreateTime = time.Now().Unix()
	ot.UpdateTime = ot.CreateTime
	// 订单号，加锁
	// mu.Lock()
	ot.RandOrder()
	if err := pdo.Db.Model(&ot).Create(&ot).Error; err != nil {
		// mu.Unlock()
		return err
	}
	// mu.Unlock()
	//记录获取订单结果
	orderParam, _ := json.Marshal(ot)
	common.ChargeLog(ot.Channel, string(orderParam))
	return nil

}

func (ot *OrderTmp) GetTmpOrderInfo() error {

	if err := pdo.Db.Model(&ot).Where("order_id = ?", ot.OrderId).First(&ot).Error; err != nil {
		return err
	}
	return nil
}

// 生成规则订单号
func (ot *OrderTmp) RandOrder() {
	mu.Lock()
	defer mu.Unlock()
	NowTime := time.Now()
	preOrder := NowTime.Format("20060102")
	midOder := CountNum()()
	suffixOrder := strconv.Itoa(int(NowTime.UnixNano() / 1e6))
	ot.OrderId = preOrder + midOder + suffixOrder[4:]

}

// 订单中间计数器
func CountNum() func() string {
	return func() string {
		if num > 10000 {
			num = 0
		}
		num++
		return fmt.Sprintf("%05d", num)
	}
}
