package main

import (
	"github.com/gin-gonic/gin"
	"github.com/zsly3n3/statisticsServer/db"
	//"github.com/zsly3n3/statisticsServer/log"
	"github.com/zsly3n3/statisticsServer/datastruct"
)

var dbHandler *db.DBHandler

func setupRouter() *gin.Engine {
	r := gin.Default()
	r.GET("/getleague", func(c *gin.Context) {
		c.JSON(200,dbHandler.GetLeague())
	})
	r.POST("/login", func(c *gin.Context) {
		name := c.Request.PostFormValue("name")
		pwd :=c.Request.PostFormValue("pwd")
		code,level:=dbHandler.Login(name,pwd)
		if isNULLError(code){
			c.JSON(200, gin.H{
				"code":code,
				"level": level,
			})
		}else{
			c.JSON(200, gin.H{
				"code":code,
			})
		}
	})
	/*插入新的关系,玩家账号与游戏id,推荐人与游戏id*/
	r.POST("/insertgidtidrid", func(c *gin.Context) {
		var body datastruct.PostGidTidRidBody
		err:=c.BindJSON(&body)
		code:=datastruct.NULLError
		if err==nil{
			code=dbHandler.InsertGidData(body.Gids,body.Tid,body.Rid)
		}else{
			code=datastruct.JsonParseFailedFromPostBody
		}
	    c.JSON(200, gin.H{
			"code":code,
	    })
	})
	return r
}

func isNULLError(code datastruct.CodeType)bool{
	 tf:=false
	 if code == datastruct.NULLError{
		tf = true
	 }
	 return tf
}

func main() {
	r := setupRouter()
	dbHandler = db.CreateDBHandler()
	r.Run(":8080")
}
