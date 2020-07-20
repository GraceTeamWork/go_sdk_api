package common

import (
	"bufio"
	"fmt"
	"os"
	"time"
)

// 登陆日志
//os 包
func LoginLog(channel string, msg string) {

	path, err := os.Getwd()
	if err != nil {
		fmt.Println("写入错误：", err)
	}
	//拼接文件路径
	path += "/log/login/" + channel
	//判断目录是否存在，不存在创建
	if !IsExist(path) {
		err := os.MkdirAll(path, os.ModePerm)
		if err != nil {
			fmt.Println("写入错误：", err)
		}
	}
	//生成文件名
	fileName := time.Now().Format("2006-01") + ".log"
	file := path + "/" + fileName
	////判断文件是否存在，不存在创建
	//if !IsExist(file) {
	//	_ ,err :=os.Create(file)
	//	if err!=nil {
	// fmt.Println("写入错误：", err) //	}
	//}
	fp, err := os.OpenFile(file, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		fmt.Println("写入错误：", err)
	}
	defer fp.Close()
	//写入数据
	data := "[" + time.Now().Format("2006-01-02 15:04:05") + "]" + msg + "\r\n"
	_, err = fp.Write([]byte(data))
	if err != nil {
		fmt.Println("写入错误：", err)
	}
}

//bufio 包

func BufLoginLog(channel string, msg string) {
	path, err := os.Getwd()
	if err != nil {
		fmt.Println("写入错误：", err)
	}
	//拼接文件路径
	path += "/log/login/" + channel
	//判断目录是否存在，不存在创建
	if !IsExist(path) {
		err := os.MkdirAll(path, os.ModePerm)
		if err != nil {
			fmt.Println("写入错误：", err)
		}
	}
	//生成文件名
	fileName := time.Now().Format("2006-01") + ".log"

	file := path + "/" + fileName

	fp, err := os.OpenFile(file, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		fmt.Println("写入错误：", err)
	}
	defer fp.Close()

	buf := bufio.NewWriter(fp)
	//写入数据
	data := "[" + time.Now().Format("2006-01-02 15:04:05") + "]" + msg + "\r\n"
	_, err = buf.Write([]byte(data))
	if err != nil {
		fmt.Println("写入错误：", err)
	}
	fmt.Printf("写入成功")
}

//判断文件或目录是否存在
func IsExist(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return false
		}
		return true
	}
	return true
}

// 下单日志
//os 包
func ChargeLog(channel string, msg string) {

	path, err := os.Getwd()
	if err != nil {
		fmt.Println("写入错误：", err)
	}
	//拼接文件路径
	path += "/log/charge/" + channel
	//判断目录是否存在，不存在创建
	if !IsExist(path) {
		err := os.MkdirAll(path, os.ModePerm)
		if err != nil {
			fmt.Println("写入错误：", err)
		}
	}
	//生成文件名
	fileName := time.Now().Format("2006-01") + ".log"
	file := path + "/" + fileName
	////判断文件是否存在，不存在创建
	//if !IsExist(file) {
	//	_ ,err :=os.Create(file)
	//	if err!=nil {
	// fmt.Println("写入错误：", err) //	}
	//}
	fp, err := os.OpenFile(file, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		fmt.Println("写入错误：", err)
	}
	defer fp.Close()
	//写入数据
	data := "[" + time.Now().Format("2006-01-02 15:04:05") + "]" + msg + "\r\n"
	_, err = fp.Write([]byte(data))
	if err != nil {
		fmt.Println("写入错误：", err)
	}
}

// 下单日志
//os 包
func CallbackLog(channel string, msg string) {

	path, err := os.Getwd()
	if err != nil {
		fmt.Println("写入错误：", err)
	}
	//拼接文件路径
	path += "/log/callback/" + channel
	//判断目录是否存在，不存在创建
	if !IsExist(path) {
		err := os.MkdirAll(path, os.ModePerm)
		if err != nil {
			fmt.Println("写入错误：", err)
		}
	}
	//生成文件名
	fileName := time.Now().Format("2006-01") + ".log"
	file := path + "/" + fileName
	fp, err := os.OpenFile(file, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		fmt.Println("写入错误：", err)
	}
	defer fp.Close()
	//写入数据
	data := "[" + time.Now().Format("2006-01-02 15:04:05") + "]" + msg + "\r\n"
	_, err = fp.Write([]byte(data))
	if err != nil {
		fmt.Println("写入错误：", err)
	}
}

// 下单日志
//os 包
func SyncLog(channel string, msg string) {

	path, err := os.Getwd()
	if err != nil {
		fmt.Println("写入错误：", err)
	}
	//拼接文件路径
	path += "/log/sync/" + channel
	//判断目录是否存在，不存在创建
	if !IsExist(path) {
		err := os.MkdirAll(path, os.ModePerm)
		if err != nil {
			fmt.Println("写入错误：", err)
		}
	}
	//生成文件名
	fileName := time.Now().Format("2006-01") + ".log"
	file := path + "/" + fileName
	fp, err := os.OpenFile(file, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		fmt.Println("写入错误：", err)
	}
	defer fp.Close()
	//写入数据
	data := "[" + time.Now().Format("2006-01-02 15:04:05") + "]" + msg + "\r\n"
	_, err = fp.Write([]byte(data))
	if err != nil {
		fmt.Println("写入错误：", err)
	}
}

// 公共日志
//os 包
func CommonLog(path string, msg string) {
	//判断目录是否存在，不存在创建
	if !IsExist(path) {
		err := os.MkdirAll(path, os.ModePerm)
		if err != nil {
			fmt.Println("写入错误：", err)
		}
	}
	//生成文件名
	fileName := time.Now().Format("2006-01") + ".log"
	file := path + "/" + fileName
	fp, err := os.OpenFile(file, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		fmt.Println("写入错误：", err)
	}
	defer fp.Close()
	//写入数据
	data := "[" + time.Now().Format("2006-01-02 15:04:05") + "]" + msg + "\r\n"
	_, err = fp.Write([]byte(data))
	if err != nil {
		fmt.Println("写入错误：", err)
	}
}
