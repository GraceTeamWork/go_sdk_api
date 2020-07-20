package pk5001

import (
	"api/common"
	"api/model"
	"encoding/json"
	"errors"
	// "fmt"
	"github.com/gin-gonic/gin"
	// "io/ioutil"
	"net/http"
)

type FengLingLogin struct {
	model.QuicksdkLogin
}

func Login(c *gin.Context) {

	var param FengLingLogin

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
	common.LoginLog(param.Channel, string(jsonParam))

	//sdk验证并获取服务器
	data, err := model.GetServer(&param)
	if err != nil {
		//记录请求错误日记
		common.LoginLog(param.Channel, string(jsonParam)+" error:"+err.Error())
		c.JSON(http.StatusOK, gin.H{"code": 400, "msg": err.Error(), "data": data})
		return
	}
	//验证成功返回服务器列表
	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "返回服务器列表", "data": data})

}

func (q *FengLingLogin) CheakParam() error {

	if q.Channel == "" {
		return errors.New("渠道不能为空")
	}
	if q.Package == 0 {
		return errors.New("包不能为空")
	}
	if q.V == "" {
		return errors.New("版本不能为空")
	}

	if q.Token == "" {
		return errors.New("token不能为空")
	}

	if q.Uid == "" {
		return errors.New("uid不能为空")
	}

	if q.ProductCode == "" {
		return errors.New("product_code不能为空")
	}

	return nil
}

func (q *FengLingLogin) CheakSdk() error {

	// url := fmt.Sprintf("http://checkuser.sdk.quicksdk.net/v2/checkUserInfo?token=%s&product_code=%suid=%dchannel_code=%s", q.Token, q.ProductCode, q.Uid, q.ChannelCode)

	// fmt.Println(url)
	// resp, err := http.Get(url)
	// if err != nil {
	// 	return err
	// }
	// defer resp.Body.Close()

	// ret, err := ioutil.ReadAll(resp.Body)
	// if err != nil {
	// 	return err
	// }

	// if string(ret) != "1" {
	// 	return errors.New("验证失败：" + string(ret))
	// }

	return nil
}

func (q *FengLingLogin) GetServer() (model.ServerReturn, error) {

	var user model.User
	user.CommonLogin = q.CommonLogin
	user.Account = q.Channel + "_" + q.Uid

	return model.DefaultGetServerList(&user, model.DefaultFilter)
}
