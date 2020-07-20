package model

type CommonLogin struct {
	Channel string `form:"channel" json:"channel"`
	Package int    `form:"package" json:"package"`
	V       string `form:"v" json:"v"`
	Ip      string
}

//quicksdk
type QuicksdkLogin struct {
	Token       string `form:"token" json:"token"`
	ProductCode string `form:"product_code" json:"product_code"`
	Uid         string `form:"uid" json:"uid"`
	ChannelCode string `form:"channel_code" json:"channel_code"`
	CommonLogin
}
