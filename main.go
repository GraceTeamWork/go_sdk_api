package main

import (
	"api/pdo"
	"api/router"
	"context"
)

func main() {

	err := pdo.InitMongo()
	if err != nil {
		panic(err)
	}
	defer pdo.Mongo.Disconnect(context.TODO())

	err = pdo.InitMysql()
	if err != nil {
		panic(err)
	}
	defer pdo.Db.Close()

	r := router.SetRouter()

	r.Run(":9003")
}
