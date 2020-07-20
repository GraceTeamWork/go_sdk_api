package router

import (
	"api/single/controller"
	"github.com/gin-gonic/gin"
)

func InitSingleRouter(r *gin.Engine) {

	fengling := r.Group("")
	{
		fengling.POST("charge", controller.Charge)
	}
}
