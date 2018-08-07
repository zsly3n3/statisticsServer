package main

import (
	"github.com/gin-gonic/gin"
	"github.com/zsly3n3/statisticsServer/db"
	//"github.com/zsly3n3/statisticsServer/log"
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
		tf,level:=dbHandler.Login(name,pwd)
		if tf{
			c.JSON(200, gin.H{
				"code":0,
				"level": level,
			})
		}else{
			c.JSON(200, gin.H{
				"code":-1,
				"message":"login falied",
			})
		}
	})
	return r
}

func main() {
	r := setupRouter()
	dbHandler = db.CreateDBHandler()
	r.Run(":8080")
}
