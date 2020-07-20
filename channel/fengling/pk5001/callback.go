package pk5001

import (
	"api/common"
	"api/conf"
	"api/model"
	"crypto/md5"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"regexp"
	"strconv"
)

type CallbackOrder struct {
	NtData  string `form:"nt_data"`
	Sign    string `form:"sign"`
	Md5Sign string `form:"md5Sign"`
	OrderInfo
}

type OrderInfo struct {
	IsTest       string  `xml:"message>is_test"`
	Channel      string  `xml:"message>channel"`
	ChannelUid   string  `xml:"message>channel_uid"`
	GameOrder    string  `xml:"message>game_order"`
	OrderNo      string  `xml:"message>order_no"`
	PayTime      string  `xml:"message>pay_time"`
	Amount       float64 `xml:"message>amount"`
	Status       string  `xml:"message>status"`
	ExtrasParams string  `xml:"message>extras_params"`
}

func Callback(c *gin.Context) {

	var param CallbackOrder
	err := c.ShouldBind(&param)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": "解析参数错误", "data": nil})
		return
	}
	jsonParam, err := json.Marshal(param)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": "解析参数错误", "data": nil})
		return
	}
	// 客户端IP
	// param.Ip = c.ClientIP()
	//记录请求日记
	common.CallbackLog("fengling", string(jsonParam))

	//获取订单
	err = model.AddOrder(&param)
	if err != nil {
		//记录请求错误日记
		common.CallbackLog("fengling", string(jsonParam)+" error:"+err.Error())
		c.JSON(http.StatusOK, gin.H{"code": 400, "msg": err.Error(), "data": nil})
		return
	}

	//验证成功返回服务器列表
	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "ok", "data": nil})
}

func (co *CallbackOrder) CheakSign() error {

	md5Key := conf.GetConfig("fengling.5001.Md5_Key").(string)

	md5Str := fmt.Sprintf("%x", md5.Sum([]byte(co.NtData+co.Sign+md5Key)))

	if co.Md5Sign != md5Str {
		return errors.New("签名错误")
	}

	return nil
}

func (co *CallbackOrder) CheakParam() error {

	if co.NtData == "" {
		return errors.New("参数错误")
	}
	if co.Sign == "" {
		return errors.New("参数错误")
	}
	if co.Md5Sign == "" {
		return errors.New("参数错误")
	}

	if err := co.Decode(); err != nil {
		return errors.New("参数错误")
	}

	if co.OrderNo == "" {
		return errors.New("参数错误")
	}

	if co.Amount <= 0 {
		return errors.New("参数错误")
	}

	if co.ExtrasParams == "" {
		return errors.New("参数错误")
	}

	return nil
}

func (co *CallbackOrder) Order() error {

	var tmpOrder model.OrderTmp
	// tmpOrder.OrderId = co.GameOrder
	tmpOrder.OrderId = "20200715000017977059"
	if err := tmpOrder.GetTmpOrderInfo(); err != nil {
		fmt.Print("无效订单:", tmpOrder.OrderId)
		return errors.New("无效订单")
	}

	//组装订单
	var order model.Order
	order.OrderNum = co.OrderNo
	order.CpOrderNum = co.GameOrder
	order.Channel = tmpOrder.Channel
	order.Package = tmpOrder.Package
	order.Server = tmpOrder.Server
	order.Money = co.Amount
	order.MoneyType = tmpOrder.MoneyType
	order.Account = tmpOrder.Account
	order.RoleId = tmpOrder.RoleId
	order.RoleLevel = tmpOrder.RoleLevel
	order.ConfId = tmpOrder.ConfId
	order.Ip = tmpOrder.Ip
	order.Mac = tmpOrder.Mac
	order.Imei = tmpOrder.Imei
	order.Idfa = tmpOrder.Idfa

	if err := order.AddOrder(); err != nil {
		return err
	}

	return nil
}

func (co *CallbackOrder) Decode() error {

	regexp := regexp.MustCompile(`\d+`)
	if regexp == nil {
		return errors.New("参数错误")
	}

	list := regexp.FindAllStringSubmatch(co.NtData, -1)

	if len(list) <= 0 {
		return errors.New("参数错误")
	}

	var data []byte
	// keys := []byte("kevv2t0jzgfgxnwse3w64akwgrzl3mxo")
	keys := []byte(conf.GetConfig("fengling.5001.Callback_Key").(string))
	lenth := len(keys)
	for k, v := range list {
		keyVar := keys[k%lenth]
		tmp, err := strconv.Atoi(v[0])
		if err != nil {
			return errors.New("参数错误")
		}
		data = append(data, byte(tmp)-(0xff&keyVar))
	}

	err := xml.Unmarshal(data, &co.OrderInfo)
	if err != nil {
		return err
	}

	return nil

}
