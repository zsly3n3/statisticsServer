package main

import (
	"github.com/gin-gonic/gin"
	"github.com/zsly3n3/statisticsServer/db"
)

var dbHandler *db.DBHandler

func setupRouter() *gin.Engine {
	r := gin.Default()
	r.GET("/getleague", func(c *gin.Context) {
		c.JSON(200,dbHandler.GetLeague())
	})
	r.POST("/login", func(c *gin.Context) {
		name := c.PostForm("name")
		pwd := c.PostForm("pwd")
		dbHandler.Login(name,pwd)
		c.JSON(200, gin.H{
			// "status":  "posted",
			// "message": message,
			// "nick":    nick,
		})
	})
	return r
}

func main() {
	r := setupRouter()
	dbHandler = db.CreateDBHandler()
	r.Run(":8080")
}
