package model

import (
	"api/conf"
)

type Zone struct {
	Id   int
	Name string
}

var zoneServer = make(map[string]interface{})

func (z *Zone) Get() interface{} {

	return conf.GetConfig("zone")

}
