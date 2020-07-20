package model

type OrderTmper interface {
	CheckParam() error
	Intert() (string, error)
}

type TmpOrder struct {
	GetOrderId int     `gorm:"getOrderId" form:"getOrderId"`
	Server     string  `gorm:"server" form:"server"`
	Package    int     `gorm:"package" form:"package"`
	RoleId     int     `gorm:"role_id" form:"role_id"`
	Money      float64 `gorm:"money" form:"money"`
	MoneyType  int     `gorm:"money_type" form:"money_type"`
	Account    string  `gorm:"account" form:"account"`
	ConfId     int     `gorm:"conf_id" form:"conf_id"`
	RoleLevel  int     `gorm:"role_level" form:"role_level"`
	Idfa       string  `gorm:"idfa" form:"idfa"`
	Imei       string  `gorm:"imei" form:"imei"`
	Mac        string  `gorm:"mac" form:"mac"`
}
