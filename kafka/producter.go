package kafka

import(
	"fmt"
	"github.com/Shopify/sarama"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
	"time"
)
//SyncProducer 模式
func Producter(c *gin.Context)  {
	fmt.Println("定义一个kafka的生产者配置变量")
	//定义一个kafka的生产者配置变量
	config := sarama.NewConfig()
	// 等待服务器所有副本都保存成功后的响应
	config.Producer.RequiredAcks =sarama.WaitForAll
	// 随机的分区类型：返回一个分区器，该分区器每次选择一个随机分区
	config.Producer.Partitioner=sarama.NewRandomPartitioner
	// 是否等待成功和失败后的响应
	config.Producer.Return.Successes=true

	//连接kafka集群和创建一个生产者
	//kafka集群地址
	var addr []string = []string{"localhost:9092"}
	producter,err := sarama.NewSyncProducer(addr,config)
	if err!=nil {
		panic(err)
	}
	//关闭连接
	defer producter.Close()

	//
	//Aproducter,err := sarama.NewAsyncProducer(addr,config)
	//defer Aproducter.AsyncClose()
	//Aproducter.Input()

	fmt.Println("构建要发送的消息")
	//构建要发送的消息
	msg := sarama.ProducerMessage{
		Topic:     "test", //主题
		//Key:       sarama.StringEncoder("key"),// 键
		Partition:int32(10), //分区
	}
	//
	fmt.Println("开始生产*************")
	for i:=1;i<100 ;i++  {

		msg.Key=sarama.StringEncoder(strconv.Itoa(i))
		//生产信息内容
		msg.Value=sarama.ByteEncoder("我是生产者，正在生成第--------->"+strconv.Itoa(i)+"个产品")
		//发送信息内容
		result,offset,err:=producter.SendMessage(&msg)
		if err!=nil {
			fmt.Println("Send message Fail")
		}
		fmt.Printf("result = %d, offset=%d\n", result, offset)

		time.Sleep(time.Second*1)
	}

	 c.JSON(http.StatusOK,gin.H{"code":200,"msg":"生成完成"})
}

//AsyncProducer 模式
func Aproducter(c *gin.Context)  {
	//定义默认的kafka配置
	config := sarama.NewConfig()
	//更改生产者配置
	//等待成功响应返回
	config.Producer.Return.Successes=true
	// 等待服务器所有副本都保存成功后的响应
	config.Producer.RequiredAcks =sarama.WaitForAll
	// 随机的分区类型：返回一个分区器，该分区器每次选择一个随机分区
	config.Producer.Partitioner=sarama.NewRandomPartitioner

	//kafka集群地址
	addr := []string{"localhost:9092"}
	//定义一个生产者
	Aproducter,err:= sarama.NewAsyncProducer(addr,config)
	if err !=nil {
		log.Fatal(err)
	}
	defer func() {
		err:=Aproducter.Close()
		if err !=nil {
			log.Fatal(err)
		}
	}()

	//构建生产消息体
	msg :=sarama.ProducerMessage{}

	msg.Topic="aProduct"

	//生产25个信息
	for i:=0;i<25 ;i++  {
		//建
		msg.Key = sarama.StringEncoder(strconv.Itoa(i))
		//值
		msg.Value = sarama.ByteEncoder("生产第"+strconv.Itoa(i)+"个产品")

		fmt.Println("生产第-------------------->"+strconv.Itoa(i)+"个产品")
		Aproducter.Input() <- &msg

		Aproducter.Successes()

		time.Sleep(time.Second*1)

	}

	c.JSON(http.StatusOK,gin.H{"code":200,"msg":"生成完成"})
}