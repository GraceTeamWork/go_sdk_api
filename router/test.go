package router

import (
	"api/test"
	"github.com/gin-gonic/gin"
)

func InitTestRouter(r *gin.Engine) {

	fengling := r.Group("")
	{
		// fengling.POST("fengling/login/5001", pk5001.Login)
		fengling.POST("test/charge", test.Charge)
		fengling.POST("test/callback", test.Callback)
	}
}
