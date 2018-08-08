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
	gameId:=&datastruct.GameId{}
	thirdPartyId:=&datastruct.ThirdPartyId{}
	referrer:=&datastruct.Referrer{}
	tg:=&datastruct.ThirdPartyId_1_n_gameId{}
	rg:=&datastruct.Referrer_1_n_gameId{}
	err:=engine.DropTables(login,league,role,gameId,thirdPartyId,referrer,tg,rg)
    errhandle(err)
	err=engine.CreateTables(login,league,role,gameId,thirdPartyId,referrer,tg,rg)
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

func (handler *DBHandler)Login(name string,pwd string)(datastruct.CodeType,int){
	tmp:=new(tmpData)
	code:= datastruct.NULLError
	if name == datastruct.NULLSTRING || pwd == datastruct.NULLSTRING{
		code = datastruct.ParamError
	}
	sql:="select level,password from login join role on role.id = login.role_id where login_name='"+name+"'"
	handler.dbEngine.Sql(sql).Get(tmp)
	level:=-1
	if tmp.Password == pwd{
	   level = tmp.Level
	}else{
	   code = datastruct.LoginFailed
	}
	return code,level
}

func (handler *DBHandler)InsertGidData(gids []string,tid string,rid string)datastruct.CodeType{
	code:= datastruct.NULLError
	if len(gids) <= 0 || (tid == datastruct.NULLSTRING && rid == datastruct.NULLSTRING){
	   code = datastruct.ParamError
	}else{
		engine:=handler.dbEngine
		session := engine.NewSession()
		defer session.Close()
		session.Begin()
		tid_id:=datastruct.NULLID
		rid_id:=datastruct.NULLID
		if tid != datastruct.NULLSTRING{
			var tmp_tid datastruct.ThirdPartyId
			tmp_tid.Identity = tid
			has, err:= session.Where("identity=?",tid).Get(&tmp_tid)
			if err != nil{
				rollback(err.Error(),session)
				return datastruct.DBSessionGetError 
			}
			if !has{
				_, err = session.Insert(&tmp_tid) 
				if err != nil{
					rollback(err.Error(),session)
					return datastruct.DBSessionInsertError
				}
			}
			tid_id = tmp_tid.Id
		}
		if rid != datastruct.NULLSTRING{
			var tmp_rid datastruct.Referrer
			tmp_rid.Identity = rid
			has, err:= session.Where("identity=?",rid).Get(&tmp_rid)
			if err != nil{
				rollback(err.Error(),session)
				return datastruct.DBSessionGetError 
			}
			if !has{
				_, err = session.Insert(&tmp_rid) 
				if err != nil{
					rollback(err.Error(),session)
					return datastruct.DBSessionInsertError
				}
			}
			rid_id = tmp_rid.Id
		}
		for _,v:=range gids{
			code = insertGid(v,tid_id,rid_id,session)
			if code != datastruct.NULLError{
			   session.Rollback()	
			   return code
			}
		}
		err:=session.Commit()
		if err != nil{
			log.Debug("session commit err:%v",err.Error())
			return datastruct.DBSessionCommitError
		}
	}
	return code
}

func insertGid(gid string,tid_id int,rid_id int,session *xorm.Session)datastruct.CodeType{
	code:= datastruct.NULLError
	if gid == datastruct.NULLSTRING{
	   return datastruct.ParamError
	}
	var tmp_gid datastruct.GameId
	tmp_gid.Identity = gid
	has, err:= session.Where("identity=?",gid).Get(&tmp_gid)
	if err != nil{
	   log.Debug("err:%v",err.Error())
	   return datastruct.DBSessionGetError 
	}
	if has{
		gid_id_str:=fmt.Sprintf("%d",tmp_gid.Id)
		if tid_id != datastruct.NULLID{
		   sql:="delete from third_party_id_1_n_game_id where g_id ="+gid_id_str
		   _,err = session.Exec(sql)
		   if err != nil{
			log.Debug("err:%v",err.Error())
			return datastruct.DBSessionExecError
		   }
		}
		if rid_id != datastruct.NULLID{
			sql:="delete from referrer_1_n_game_id where g_id ="+gid_id_str
			_,err = session.Exec(sql)
			if err != nil{
			 log.Debug("err:%v",err.Error())
			 return datastruct.DBSessionExecError
			}
		}
	}else{
		_, err = session.Insert(&tmp_gid) 
		if err != nil{
			log.Debug("err:%v",err.Error())
			return datastruct.DBSessionInsertError
		}
	}
	gid_id:=tmp_gid.Id
	if tid_id != datastruct.NULLID{
	   var t_1_n_g datastruct.ThirdPartyId_1_n_gameId
	   t_1_n_g.GId = gid_id
	   t_1_n_g.TId = tid_id
	   _, err = session.Insert(&t_1_n_g) 
	   if err != nil{
		   log.Debug("err:%v",err.Error())
		   return datastruct.DBSessionInsertError
	   }
	}
	if rid_id != datastruct.NULLID{
		var r_1_n_g datastruct.Referrer_1_n_gameId
		r_1_n_g.GId = gid_id
		r_1_n_g.RId = rid_id
		_, err = session.Insert(&r_1_n_g) 
		if err != nil{
			log.Debug("err:%v",err.Error())
			return datastruct.DBSessionInsertError
		}
	}
	return code
}

func rollback(err_str string,session *xorm.Session){
	 log.Debug("will rollback,err_str:%v",err_str)
	 session.Rollback()
}