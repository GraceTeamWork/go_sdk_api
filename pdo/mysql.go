package pdo

import (
	"api/conf"
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var Db *gorm.DB

// 初始化数据库连接
func InitMysql() (err error) {
	dsn := getDsn()
	Db, err = gorm.Open("mysql", dsn) //链接数据库
	if err != nil {
		return err
	}
	err = Db.DB().Ping()
	Db.SingularTable(true) //gorm默认把表名都加上 "s"，关掉默认
	//设置表前缀
	gorm.DefaultTableNameHandler = func(db *gorm.DB, defaultTableName string) string {
		return "xxwd_" + defaultTableName
	}
	return nil
}

func getDsn() string {
	username := conf.GetConfig("mysql.username")
	password := conf.GetConfig("mysql.password")
	host := conf.GetConfig("mysql.host")
	port := conf.GetConfig("mysql.port")
	dbname := conf.GetConfig("mysql.dbname")
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8", username, password, host, port, dbname)
}
