package router

import (
	"api/channel/fengling/pk5001"
	"github.com/gin-gonic/gin"
)

func InitChannelFengLingRouter(r *gin.Engine) {

	fengling := r.Group("/channel/")
	{
		fengling.POST("fengling/login/5001", pk5001.Login)
		fengling.POST("fengling/charge/5001", pk5001.Charge)
		fengling.POST("fengling/callback/5001", pk5001.Callback)
	}
}
