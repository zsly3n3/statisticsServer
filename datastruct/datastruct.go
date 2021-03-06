package datastruct

import (
	"time"
)

const NULLSTRING = ""
const NULLID = -1

type CodeType int //错误码
const (
	NULLError CodeType = iota //无错误
	ParamError//参数错误,数据为空或者类型不对等
	LoginFailed//登录失败,如无此账号或者密码错误等
	JsonParseFailedFromPostBody//来自post请求中的Body解析json失败
	DBSessionGetError//xorm事务中Get方法执行出错
	DBSessionExecError//xorm事务中Exec方法执行出错
	DBSessionInsertError//xorm事务中Insert方法执行出错
	DBSessionCommitError//xorm事务中Commit方法执行出错
	DBSessionUpdateError//xorm事务中Update方法执行出错
)

type Role struct {
	Id    int       `xorm:"not null pk autoincr INT(11)"`
	Level int       `xorm:"INT(11) not null"`  //权限等级
    Desc  string    `xorm:"VARCHAR(32) not null"` //权限名称
}

type League struct {
    Id   int       `xorm:"not null pk autoincr INT(11)"`
	Name string    `xorm:"VARCHAR(32) not null"`
}

type Login struct {
    Id        int       `xorm:"not null pk autoincr INT(11)"`
	LoginName string    `xorm:"VARCHAR(64) not null"`
	Password  string    `xorm:"VARCHAR(128) not null"`
	RoleId    int       `xorm:"INT(11) not null"` //权限id
	CreatedAt time.Time `xorm:"created"`
}

/*游戏Id*/
type GameId struct {
    Id       int       `xorm:"not null pk autoincr INT(11)"`
	Identity string    `xorm:"VARCHAR(256) not null"`
}

/*玩家账号,比如微信账号*/
type ThirdPartyId struct {
    Id       int       `xorm:"not null pk autoincr INT(11)"`
	Identity string    `xorm:"VARCHAR(256) not null"`
}

/*推荐人*/
type Referrer struct {
	Id   int       `xorm:"not null pk autoincr INT(11)"`
	Identity string    `xorm:"VARCHAR(256) not null"`
}

/*玩家账号与游戏Id,一对多关系*/
type ThirdPartyId_1_n_gameId struct {
	Id  int   `xorm:"not null pk autoincr INT(11)"`
	TId int   `xorm:"INT(11) not null"`
	GId int   `xorm:"INT(11) not null"`
}

/*玩家账号与推荐人,一对一关系*/
type ThirdPartyId_1_1_referrerId struct {
	Id  int   `xorm:"not null pk autoincr INT(11)"`
	TId int   `xorm:"INT(11) not null"`
	RId int   `xorm:"INT(11) not null"`
}

/*玩家账号与利率,一对一关系*/
type ThirdPartyId_1_1_accrual struct {
	Id  int   `xorm:"not null pk autoincr INT(11)"`
	TId int   `xorm:"INT(11) not null"`
	Csl float64  `xorm:"Numeric not null"`
	Bxfl float64   `xorm:"Numeric not null"`
	Tjrfbxl float64  `xorm:"Numeric"`
}

type PostGidTidRidBody struct {
	Tid string `json:"tid"`
	Rid string `json:"rid"`
	Gids []string `json:"gids"`
	Csl float64 `json:"csl"`//抽水率
	Bxfl float64 `json:"bxfl"`//保险返率
	Tjrfbxl float64 `json:"tjrfbxl"`//推荐人返保险返率
}

type PostGidBody struct {
	Gids []string `json:"gids"`
}

type PostTidRidBody struct {
	Tid string `json:"tid"`
	Tid_id int `json:"tid_id"`
	Rid string `json:"rid"`
	Rid_id int `json:"rid_id"`
	Csl float64 `json:"csl"`//抽水率
	Bxfl float64 `json:"bxfl"`//保险返率
	Tjrfbxl float64 `json:"tjrfbxl"`//推荐人返保险返率
}

type LoginBody struct {
	LoginName string `json:"name"`
	Pwd string `json:"pwd"`
}
