package router

import (
	"api/kafka"
	"github.com/gin-gonic/gin"
)

func InitKafkaRouter(r *gin.Engine)  {

	kafkaGroup :=r.Group("/kafka/")
	//kafkaGroup :=r.Use()
	{
		kafkaGroup.POST("producter", kafka.Producter)//生产者
		kafkaGroup.POST("aproducter", kafka.Aproducter)//生产者
		kafkaGroup.POST("comsumer", kafka.Consumer)//消费者
		kafkaGroup.POST("acomsumer", kafka.Acomsumer)//消费者
	}
}
