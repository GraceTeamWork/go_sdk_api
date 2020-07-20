package router

import (
	"api/elasticsearch"
	// "fmt"
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	// "io"
	// "os"
	// "time"
)

func SetRouter() *gin.Engine {

	//gin.DisableConsoleColor()
	// 创建记录日志的文件
	//path :="/log/gin"
	//if !common.IsExist(path) {
	//	err := os.MkdirAll(path, os.ModePerm)
	//	if err != nil {
	//		fmt.Println("写入错误：", err)
	//	}
	//}
	////生成文件名
	//fileName := time.Now().Format("2006-01") + ".log"
	//file := path + "/" + fileName
	//fp, _ := os.OpenFile(file, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	//gin.DefaultWriter = io.MultiWriter(fp)

	// fileName := time.Now().Format("2006-01") + ".log"
	// f, err := os.Create("./log/gin/" + fileName)
	// if err != nil {
	// 	fmt.Println("创建文件错误：", err)
	// }
	// gin.DefaultWriter = io.MultiWriter(f)

	r := gin.Default()

	//
	pprof.Register(r)

	// 登陆、订单路由
	// InitChannelRouter(r)

	InitChannelFengLingRouter(r)

	//测试路由
	InitTestRouter(r)

	//kafka路由
	InitKafkaRouter(r)

	// 单服
	InitSingleRouter(r)

	elasticsearchGroup := r.Group("/elasticsearch/")
	{
		elasticsearchGroup.POST("/es/get", elasticsearch.Get)
	}

	return r
}
