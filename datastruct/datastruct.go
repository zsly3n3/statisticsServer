package datastruct

import (
	"time"
)

const NULLSTRING = ""
const NULLID = -1

type Role struct {
	Id    int       `xorm:"not null pk autoincr INT(11)"`
	Level int    `xorm:"INT(11) not null"`  //权限等级
    Desc  string    `xorm:"VARCHAR(32) not null"` //权限名称
}

type League struct {
    Id       int       `xorm:"not null pk autoincr INT(11)"`
	Name string    `xorm:"VARCHAR(32) not null"`
}

type Login struct {
    Id       int       `xorm:"not null pk autoincr INT(11)"`
	LoginName string    `xorm:"VARCHAR(64) not null"`
	Password string    `xorm:"VARCHAR(128) not null"`
	RoleId int    `xorm:"INT(11) not null"` //权限id
	CreatedAt time.Time `xorm:"created"`
}






