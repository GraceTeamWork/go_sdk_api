package kafka

import (
	"fmt"
	"github.com/Shopify/sarama"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
)
var (
	wg  sync.WaitGroup
)

func Consumer(c *gin.Context)  {

	//定义一个消费者
	consumer,err:= sarama.NewConsumer([]string{"localhost:9092"},nil)
	if err!=nil {
		panic(err)
	}
	//获取topic为test 的所有分区
	partitionList, err :=consumer.Partitions("aProduct")
	if err != nil {
		panic(err)
	}
	fmt.Println(partitionList)
	for partition:=range partitionList {
		//获取一个分区的所有数据
		pc, err := consumer.ConsumePartition("test", int32(partition), sarama.OffsetNewest)
		if err != nil {
			panic(err)
		}
		defer pc.AsyncClose()

		wg.Add(1)
		go func(sarama.PartitionConsumer) {
			defer wg.Done()
			//获取一个分区每条数据
			//Messages()该方法返回一个消费消息类型的只读通道，由代理产生
			for msg := range pc.Messages() {
				fmt.Printf("%s---Partition:%d, Offset:%d, Key:%s, Value:%s\n", msg.Topic,msg.Partition, msg.Offset, string(msg.Key), string(msg.Value))
			}
		}(pc)
	}
	wg.Wait()
	consumer.Close()

	c.JSON(http.StatusOK,gin.H{"code":200,"msg":"消费完成"})

}

func Acomsumer(c *gin.Context)  {
	consumer, err := sarama.NewConsumer([]string{"localhost:9092"}, nil)
	if err != nil {
		panic(err)
	}

	defer func() {
		if err := consumer.Close(); err != nil {
			log.Fatalln(err)
		}
	}()

	partitionConsumer, err := consumer.ConsumePartition("aProduct", 0, sarama.OffsetNewest)
	if err != nil {
		panic(err)
	}

	defer func() {
		if err := partitionConsumer.Close(); err != nil {
			log.Fatalln(err)
		}
	}()

	// Trap SIGINT to trigger a shutdown.
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)

	consumed := 0
ConsumerLoop:
	for {
		select {
		case msg := <-partitionConsumer.Messages():
			log.Printf("Consumed message offset %d\n", msg.Offset)
			consumed++
		case <-signals:
			break ConsumerLoop
		}
	}

	log.Printf("Consumed: %d\n", consumed)
}
