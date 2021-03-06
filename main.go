package main

import (
	"github.com/gin-gonic/gin"
	"github.com/zsly3n3/statisticsServer/db"
	//"github.com/zsly3n3/statisticsServer/log"
	"github.com/zsly3n3/statisticsServer/datastruct"
)

var dbHandler *db.DBHandler

func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Add("Access-Control-Allow-Origin", "*")
		c.Next()
	}
}
func setupRouter() *gin.Engine {
	r := gin.Default()
	r.Use(Cors())
	r.GET("/getleague", func(c *gin.Context) {
		c.JSON(200,dbHandler.GetLeague())
	})
	r.POST("/login", func(c *gin.Context) {
		var body datastruct.LoginBody
		err:=c.BindJSON(&body)
		code:=datastruct.NULLError
		if err == nil {
		   code,level:=dbHandler.Login(body.LoginName,body.Pwd)
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
		}else{
		   code=datastruct.JsonParseFailedFromPostBody
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
			code=dbHandler.InsertGidData(&body)
		}else{
			code=datastruct.JsonParseFailedFromPostBody
		}
	    c.JSON(200, gin.H{
			"code":code,
	    })
	})

	r.GET("/getGTRWithGName/:name", func(c *gin.Context) {
		g_name := c.Params.ByName("name")
		rs:=dbHandler.QueryWithGid(g_name)
		c.JSON(200, gin.H{
			"data":rs,
		})
	})

	r.GET("/getGTRWithTName/:name", func(c *gin.Context) {
		g_name := c.Params.ByName("name")
		rs:=dbHandler.QueryWithTid(g_name)
		c.JSON(200, gin.H{
			"data":rs,
		})
	})

	r.GET("/getGTRWithRName/:name", func(c *gin.Context) {
		g_name := c.Params.ByName("name")
		rs:=dbHandler.QueryWithRid(g_name)
		c.JSON(200, gin.H{
			"data":rs,
		})
	})

	r.GET("/getTRWithGName/:name", func(c *gin.Context) {
		g_name := c.Params.ByName("name")
		rs:=dbHandler.GetTR(g_name)
		c.JSON(200, gin.H{
			"data":rs,
		})
	})

	r.POST("/getTRWithGNames", func(c *gin.Context) {
		var body datastruct.PostGidBody
		err:=c.BindJSON(&body)
		rs:=make([]*datastruct.PostTidRidBody,0,len(body.Gids))
		if err==nil{
		   for _,v := range body.Gids{
			   tmp:=dbHandler.GetTR(v)
			   rs=append(rs,tmp)
		   }
		}
	    c.JSON(200, gin.H{
			"data":rs,
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
	r.Run("0.0.0.0:8181")
}
