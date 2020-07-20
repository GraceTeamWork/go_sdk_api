package test

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"strconv"
	"sync"
	"time"
)

var wg = sync.WaitGroup{}

func Charge(c *gin.Context) {

	t1 := time.Now()
	for i := 0; i < 100; i++ {

		wg.Add(1)
		go httpCharge(i)
	}

	wg.Wait()
	t2 := time.Now()
	fmt.Println("ok")
	fmt.Println(t1)
	fmt.Println(t2)
}

func httpCharge(j int) {

	defer wg.Done()

	var param = url.Values{}

	url := "http://127.0.0.1:9003/channel/fengling/charge/5001"

	param.Set("getOrderId", "1")

	rand.Seed(time.Now().UnixNano())

	for i := 0; i <= 100; i++ {
		rands := rand.Intn(99999) + 10000
		param.Set("getOrderId", "1")
		param.Set("server", "ly"+"_"+strconv.Itoa(j)+strconv.Itoa(i))
		param.Set("channel", "fengling")
		param.Set("package", strconv.Itoa(i))
		param.Set("role_id", strconv.Itoa(rands))
		param.Set("account", "fengling_"+"_"+strconv.Itoa(j)+strconv.Itoa(rands))
		param.Set("money", strconv.Itoa(rand.Intn(99999)+1))
		param.Set("money_type", strconv.Itoa(rand.Intn(7)+1))
		param.Set("conf_id", strconv.Itoa(rand.Intn(100)+1))
		param.Set("role_level", strconv.Itoa(rand.Intn(100)+10))

		resp, err := http.PostForm(url, param)
		if err != nil {
			log.Fatal(err)
			return
		}
		resp.Body.Close()
	}
}
