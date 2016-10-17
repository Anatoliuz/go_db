package main

import "github.com/gin-gonic/gin"

func (db *Resource) commonClear(context *gin.Context) {
	db.Map.Exec("TRUNCATE TABLE post")
	db.Map.Exec("TRUNCATE TABLE user")
	db.Map.Exec("TRUNCATE TABLE forum")
	db.Map.Exec("TRUNCATE TABLE follow")
	db.Map.Exec("TRUNCATE TABLE thread")
	db.Map.Exec("TRUNCATE TABLE subscription")
	context.JSON(200, gin.H{"code": 0, "response": "OK"})
}

func (db *Resource) commonStatus(context *gin.Context) {
	users, _ := db.Map.SelectInt("SELECT count(id) FROM user")
	threads, _ := db.Map.SelectInt("SELECT count(id) FROM thread")
	forums, _ := db.Map.SelectInt("SELECT count(id) FROM forum")
	posts, _ := db.Map.SelectInt("SELECT count(id) FROM post")
	context.JSON(200, gin.H{"code": 0, "response": gin.H{"user": users, "thread": threads, "forum": forums, "post": posts}})
}
