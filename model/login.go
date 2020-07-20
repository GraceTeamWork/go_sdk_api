package model

import (
	"sync"
)

var wg sync.WaitGroup

var Lock = sync.RWMutex{}

type Loginer interface {
	CheakParam() error
	CheakSdk() error
	GetServer() (ServerReturn, error)
}

type ServerReturn struct {
	Zone           interface{} `json:"zone"`
	Server         []*Server   `json:"server"`
	RecentlyServer []string    `json:"recently_server"`
	User           User        `json:"user"`
}

func GetServer(channel Loginer) (returnData ServerReturn, err error) {

	if err = channel.CheakParam(); err != nil {
		return
	}

	if err = channel.CheakSdk(); err != nil {
		return
	}

	return channel.GetServer()
}

func DefaultGetServerList(user *User, DefaultFilters func(User, []*Server) (data []*Server)) (returnData ServerReturn, err error) {
	// 服务器列表
	var server Server
	//var ch =make(chan struct{})
	serverLists, err := server.GetServerList(user)
	returnData.Server =serverLists
	if err != nil {
		return returnData, err
	}
	// 过滤服务器
	Lock.Lock()
	defer Lock.Unlock()
	var ok = make(chan struct{}, 3)
	wg.Add(3)
	go func() {
		defer wg.Done()
		returnData.Server = DefaultFilters(*user, serverLists)
		ok <- struct{}{}
	}()
	// 获取大区
	var zoneList map[string]interface{}
	go func() {
		defer wg.Done()
		var zone Zone
		zoneList = zone.Get().(map[string]interface{})
		ok <- struct{}{}
	}()
	// 角色信息
	go func() {
		defer wg.Done()
		returnData.User = *user
		// 最近登陆
		userInfo := user.GetUserInfo()
		returnData.RecentlyServer = GetMapKey(&userInfo)
		ok <- struct{}{}
	}()
	wg.Wait()
	close(ok)
	for {
		_, isOk := <-ok
		if !isOk {
			break
		}
	}
	for k, _ := range zoneServer {
		if zoneList[k] != nil {
			zoneServer[k] = zoneList[k]
		} else {
			zoneServer[k] = "特殊区"
		}
	}
	returnData.Zone = zoneServer

	return returnData, nil

}

func GetMapKey(m *map[string]interface{}) (keys []string) {
	for k, _ := range *m {
		keys = append(keys, k)
	}
	return
}
