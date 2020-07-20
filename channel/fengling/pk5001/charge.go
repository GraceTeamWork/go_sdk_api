package pk5001

import (
	"api/common"
	"api/model"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Charge(c *gin.Context) {

	var param model.OrderTmp
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
	param.Ip = c.ClientIP()
	//记录请求日记
	common.ChargeLog(param.Channel, string(jsonParam))

	//获取订单
	err = model.GetOrderNum(&param)
	if err != nil {
		//记录请求错误日记
		common.ChargeLog(param.Channel, string(jsonParam)+" error:"+err.Error())
		c.JSON(http.StatusOK, gin.H{"code": 400, "msg": err.Error(), "data": nil})
		return
	}

	//验证成功返回服务器列表
	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "ok", "data": param.OrderId})

}
