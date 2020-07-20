package controller

import (
	"api/single/model"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Charge(c *gin.Context) {
	var param model.Order
	if err := c.ShouldBind(&param); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": "解析参数错误", "data": nil})
		return
	}
	//获取订单
	if err := param.AddOrder(); err != nil {
		//记录请求错误日记
		// common.CallbackLog("fengling", string(jsonParam)+" error:"+err.Error())
		c.JSON(http.StatusOK, gin.H{"code": 400, "msg": err.Error(), "data": nil})
		return
	}

	//验证成功返回服务器列表
	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "success", "data": nil})
}
