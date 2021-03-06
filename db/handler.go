package db

import (
	"fmt"
	"github.com/go-xorm/xorm"
	"github.com/zsly3n3/statisticsServer/log"
	"github.com/zsly3n3/statisticsServer/datastruct"
	_ "github.com/go-sql-driver/mysql"
)

const DB_IP = "14.29.123.151:3306"
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
	tr:=&datastruct.ThirdPartyId_1_1_referrerId{}
	ta:=&datastruct.ThirdPartyId_1_1_accrual{} 
	err:=engine.DropTables(login,league,role,gameId,thirdPartyId,referrer,tg,tr,ta)
    errhandle(err)
	err=engine.CreateTables(login,league,role,gameId,thirdPartyId,referrer,tg,tr,ta)
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
		Name:"A",
	}
	b:= datastruct.League{
		Name:"B",
	}
	c:= datastruct.League{
		Name:"C",
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

func (handler *DBHandler)InsertGidData(body *datastruct.PostGidTidRidBody)datastruct.CodeType{
	gids:=body.Gids
	tid:=body.Tid
	rid:=body.Rid
	csl:=body.Csl
	bxfl:=body.Bxfl
	tjrfbxl:=body.Tjrfbxl
	code:= datastruct.NULLError
	if len(gids) <= 0 || tid == datastruct.NULLSTRING{
	   code = datastruct.ParamError
	}else{
		engine:=handler.dbEngine
		session := engine.NewSession()
		defer session.Close()
		session.Begin()
		tid_id:=datastruct.NULLID
		rid_id:=datastruct.NULLID
        var sql string
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
			tid_id_str:=fmt.Sprintf("%d",tid_id)
		    if tid_id != datastruct.NULLID{
		      sql="delete from third_party_id_1_n_game_id where t_id ="+tid_id_str
		      _,err = session.Exec(sql)
		      if err != nil{
				rollback(err.Error(),session)
			    return datastruct.DBSessionExecError
			  }

			  sql="delete from third_party_id_1_1_referrer_id where t_id ="+tid_id_str
			  _,err = session.Exec(sql)
			  if err != nil{
				rollback(err.Error(),session)
				return datastruct.DBSessionExecError
			  }
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
			if rid_id != datastruct.NULLID{
		        var t_1_1_r datastruct.ThirdPartyId_1_1_referrerId
		        t_1_1_r.TId = tid_id
		        t_1_1_r.RId = rid_id
		        _, err = session.Insert(&t_1_1_r)
		        if err != nil{
			      rollback(err.Error(),session)
			      return datastruct.DBSessionInsertError
				}
	        }
		}
		var t_1_1_a datastruct.ThirdPartyId_1_1_accrual
		has, err= session.Where("t_id=?",tid_id).Get(&t_1_1_a)
		if err != nil{
		  rollback(err.Error(),session)
		  return datastruct.DBSessionGetError 
		}
		t_1_1_a.TId = tid_id
		t_1_1_a.Csl = csl
		t_1_1_a.Bxfl = bxfl
		t_1_1_a.Tjrfbxl = tjrfbxl
		if !has{
			_, err = session.Insert(&t_1_1_a) 
			if err != nil{
				rollback(err.Error(),session)
				return datastruct.DBSessionInsertError
			}
		}else{
			_, err = session.Where("t_id=?",tid_id).Update(&t_1_1_a)
			if err != nil{
				rollback(err.Error(),session)
				return datastruct.DBSessionUpdateError
			}
		}
		for _,v:=range gids{
			code = insertGid(v,tid_id,rid_id,session)
			if code != datastruct.NULLError{
			   session.Rollback()	
			   return code
			}
		}
		err=session.Commit()
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
	if !has{
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
	return code
}

func rollback(err_str string,session *xorm.Session){
	 log.Debug("will rollback,err_str:%v",err_str)
	 session.Rollback()
}

func (handler *DBHandler)QueryWithGid(name string)*datastruct.PostGidTidRidBody{
	engine:=handler.dbEngine
	var gid datastruct.GameId
	has, _:= engine.Where("identity=?",name).Get(&gid)
	if has{
	   sql:="select third_party_id.id,third_party_id.identity from game_id join third_party_id_1_n_game_id on third_party_id_1_n_game_id.g_id = game_id.id join third_party_id on third_party_id_1_n_game_id.t_id = third_party_id.id where game_id.id="+fmt.Sprintf("%d",gid.Id)
	   var tid datastruct.ThirdPartyId
	   has_tid,_:=engine.Sql(sql).Get(&tid)
	   if has_tid{
		  return getGidTidRidData(&tid,engine,datastruct.NULLSTRING)
	   }  
	   return nil
	}
	return nil
}
func (handler *DBHandler)QueryWithTid(t_name string)*datastruct.PostGidTidRidBody{
	engine:=handler.dbEngine
	var tid datastruct.ThirdPartyId
	has, _:= engine.Where("identity=?",t_name).Get(&tid)
	if has{
	   return getGidTidRidData(&tid,engine,datastruct.NULLSTRING)
	}
	return nil
}

func (handler *DBHandler)QueryWithRid(r_name string)*datastruct.PostGidTidRidBody{
	engine:=handler.dbEngine
	var rid datastruct.Referrer
	has, _:= engine.Where("identity=?",r_name).Get(&rid)
	if has{
	   sql:="select third_party_id.id,third_party_id.identity from third_party_id join third_party_id_1_1_referrer_id on third_party_id_1_1_referrer_id.t_id = third_party_id.id where r_id="+fmt.Sprintf("%d",rid.Id)
	   var tid datastruct.ThirdPartyId
	   has_tid,_:=engine.Sql(sql).Get(&tid)
	   if has_tid{
		  return getGidTidRidData(&tid,engine,rid.Identity)
	   }  
	   return nil
	}
	return nil
}

type tmpGIdentity struct {
	Identity string
}
func getGidTidRidData(tid *datastruct.ThirdPartyId,engine *xorm.Engine,r_name string)*datastruct.PostGidTidRidBody{
	   rs:=new(datastruct.PostGidTidRidBody)
	   rs.Tid = tid.Identity
	   var identitys []tmpGIdentity
	   sql:="select game_id.identity from game_id join third_party_id_1_n_game_id on third_party_id_1_n_game_id.g_id = game_id.id where t_id="+fmt.Sprintf("%d",tid.Id)
	   engine.Sql(sql).Find(&identitys)
	   gids:=make([]string,0,len(identitys))
	   for _,v :=range identitys {
		 gids = append(gids,v.Identity)
	   }
	   rs.Gids = gids
       if r_name == datastruct.NULLSTRING{
		  var referrer_data datastruct.Referrer
		  sql="select referrer.id,referrer.identity from referrer join third_party_id_1_1_referrer_id on third_party_id_1_1_referrer_id.r_id = referrer.id where t_id="+fmt.Sprintf("%d",tid.Id)
		  has, _:=engine.Sql(sql).Get(&referrer_data)
		  if has {
		    rs.Rid = referrer_data.Identity
		  }
	   }else{
		   rs.Rid = r_name
	   }
	   var t_1_1_a datastruct.ThirdPartyId_1_1_accrual
	   engine.Where("t_id=?",tid.Id).Get(&t_1_1_a)
	   rs.Csl = t_1_1_a.Csl
	   rs.Bxfl = t_1_1_a.Bxfl
	   rs.Tjrfbxl = t_1_1_a.Tjrfbxl
	   return rs
}

func (handler *DBHandler)GetTR(name string)*datastruct.PostTidRidBody{
	engine:=handler.dbEngine
	var gid datastruct.GameId
	has, _:= engine.Where("identity=?",name).Get(&gid)
	if has{
	   sql:="select third_party_id.id,third_party_id.identity from game_id join third_party_id_1_n_game_id on third_party_id_1_n_game_id.g_id = game_id.id join third_party_id on third_party_id_1_n_game_id.t_id = third_party_id.id where game_id.id="+fmt.Sprintf("%d",gid.Id)
	   var tid datastruct.ThirdPartyId
	   has_tid,_:=engine.Sql(sql).Get(&tid)
	   if has_tid{
		  return getTidRidData(&tid,engine)
	   }
	   return nil
	}
	return nil
}

func getTidRidData(tid *datastruct.ThirdPartyId,engine *xorm.Engine)*datastruct.PostTidRidBody{
	rs:=new(datastruct.PostTidRidBody)
	rs.Tid = tid.Identity
	rs.Tid_id = tid.Id
	var referrer_data datastruct.Referrer
	sql:="select referrer.id,referrer.identity from referrer join third_party_id_1_1_referrer_id on third_party_id_1_1_referrer_id.r_id = referrer.id where t_id="+fmt.Sprintf("%d",tid.Id)
	has, _:=engine.Sql(sql).Get(&referrer_data)
	if has {
		 rs.Rid = referrer_data.Identity
		 rs.Rid_id = referrer_data.Id
	}
	var t_1_1_a datastruct.ThirdPartyId_1_1_accrual
	engine.Where("t_id=?",tid.Id).Get(&t_1_1_a)
	rs.Csl = t_1_1_a.Csl
	rs.Bxfl = t_1_1_a.Bxfl
	rs.Tjrfbxl = t_1_1_a.Tjrfbxl
	return rs
}