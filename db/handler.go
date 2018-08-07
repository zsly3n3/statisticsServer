package db

import (
	"fmt"
	"github.com/go-xorm/xorm"
	"github.com/zsly3n3/statisticsServer/log"
	"github.com/zsly3n3/statisticsServer/datastruct"
	_ "github.com/go-sql-driver/mysql"
)

const DB_IP = "47.105.45.83:3306"
const DB_Name = "statistics"
const DB_UserName = "root"
const DB_Pwd = "zsly3n@S"

type DBHandler struct {
	 dbEngine *xorm.Engine
}

func CreateDBHandler()*DBHandler{
	dbHandler:=new(DBHandler)
    dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8",DB_UserName,DB_Pwd,DB_IP,DB_Name)
	engine, err:= xorm.NewEngine("mysql", dsn)
	errhandle(err)
	err=engine.Ping()
	errhandle(err)
	//日志打印SQL
    engine.ShowSQL(true)
	//设置连接池的空闲数大小
	engine.SetMaxIdleConns(1)
    resetDB(engine)
    initData(engine)
	dbHandler.dbEngine = engine
    return dbHandler
}

func errhandle(err error){
	if err != nil {
		log.Fatal("db error is %v", err.Error())
	}
}

func resetDB(engine *xorm.Engine){
	login:=&datastruct.Login{}
    league:=&datastruct.League{}
    role:=&datastruct.Role{}
	err:=engine.DropTables(login,league,role)
    errhandle(err)
	err=engine.CreateTables(login,league,role)
    errhandle(err)
}

func initData(engine *xorm.Engine){
	league:=createLeagueData()
    _, err := engine.Insert(&league)
	errhandle(err)
	role:=createRoleData()
    _, err = engine.Insert(&role)
	errhandle(err)
	login:=createLoginData(getLevelID(engine,role))
    _, err = engine.Insert(&login)
    errhandle(err)
}

func getLevelID(engine *xorm.Engine,roles []datastruct.Role)(int,int){
    LevelID:=make([]int,0,len(roles))
	for _,v := range roles{
		var role datastruct.Role
		has,err:=engine.Where("level =?", v.Level).Get(&role)
		if !has || err !=nil{
		   log.Fatal("db query RoleID error")
		   break
		}
		LevelID = append(LevelID,role.Id)
	}
	return LevelID[0],LevelID[1]
}

func createLeagueData()[]datastruct.League{
	a:= datastruct.League{
		Name:"a",
	}
	b:= datastruct.League{
		Name:"b",
	}
	c:= datastruct.League{
		Name:"c",
	}
	return []datastruct.League{a,b,c}
}

func createRoleData()[]datastruct.Role{
	admin:= datastruct.Role{
		Level:0,
		Desc:"admin",
	}
	guest:= datastruct.Role{
		Level:1,
		Desc:"guest",
	}
	return []datastruct.Role{admin,guest}
}

func createLoginData(adminLevelID int,guestLevelID int)[]datastruct.Login{
	admin:= datastruct.Login{
		LoginName:"admin",
		Password:"123@s678",
		RoleId:adminLevelID,
	}
	guest:= datastruct.Login{
		LoginName:"guest",
		Password:"1234s6",
		RoleId:guestLevelID,
	}
	return []datastruct.Login{admin,guest}
}

func (handler *DBHandler)GetLeague()[]datastruct.League{
	var league []datastruct.League
	handler.dbEngine.Find(&league)
	return league
}

type tmpData struct {
	Level int
	Password string
}

func (handler *DBHandler)Login(name string,pwd string)(bool,int){
	tmp:=new(tmpData)
	sql:="select level,password from login join role on role.id = login.role_id where login_name='"+name+"'"
	handler.dbEngine.Sql(sql).Get(tmp)
	tf:=false
	level:=-1
	if tmp.Password == pwd{
	   tf = true
	   level = tmp.Level
	}
	return tf,level
}