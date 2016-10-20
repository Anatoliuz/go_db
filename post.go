package main

import (
	"strconv"
	"github.com/gin-gonic/gin"
)

// Post is a model
type Post struct {
	Date          string `db:"date" json:"date"`
	Dislikes      int    `db:"dislikes" json:"dislikes"`
	Forum         string `db:"forum" json:"forum"`
	ID            int    `db:"id" json:"id"`
	IsApproved    bool   `db:"isApproved" json:"isApproved"`
	IsDeleted     bool   `db:"isDeleted" json:"isDeleted"`
	IsEdited      bool   `db:"isEdited" json:"isEdited"`
	IsHighlighted bool   `db:"isHighlighted" json:"isHighlighted"`
	IsSpam        bool   `db:"isSpam" json:"isSpam"`
	Likes         int    `db:"likes" json:"likes"`
	Message       string `db:"message" json:"message"`
	Parent        *int    `db:"parent" json:"parent"`
	Points        int    `db:"points" json:"points"`
	Thread        int    `db:"thread" json:"thread"`
	User          string `db:"user" json:"user"`
	FirstPath     int    `db:"first_path" json:"first_path"`
	LastPath      string `db:"last_path" json:"last_path"`
}

func (db *Resource) postWithID(id int) gin.H {
	var post Post

	if err := db.Map.SelectOne(&post, "SELECT * FROM post WHERE id = ?", id); err == nil {

		return gin.H{"date": post.Date, "dislikes": post.Dislikes, "forum": post.Forum, "id": post.ID, "isApproved": post.IsApproved, "isDeleted": post.IsDeleted, "isEdited": post.IsEdited, "isHighlighted": post.IsHighlighted, "isSpam": post.IsSpam, "likes": post.Likes, "message": post.Message, "parent": post.Parent, "points": post.Points, "thread": post.Thread, "user": post.User, "first_path": 0, "last_path": ""}
	}else{
		//println(err.Error())
}
	return nil
}

func (db *Resource) postCreate(context *gin.Context){
	var post Post
	context.BindJSON(&post)
	postMap, _ := db.Map.Exec("INSERT INTO post (date, forum, isApproved, isDeleted, isEdited, isHighlighted, isSpam, message, parent, thread, user) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)", post.Date, post.Forum, post.IsApproved, post.IsDeleted, post.IsEdited, post.IsHighlighted, post.IsSpam, post.Message, post.Parent, post.Thread, post.User)
	id, _ := postMap.LastInsertId()

	var count int
	db.Map.SelectOne(&count, "SELECT COUNT(*) FROM user_forum where user=? AND forum=?", post.User, post.Forum)
	if count == 0{
		db.Map.Exec("INSERT INTO user_forum (user, forum) VALUES(?,?)", post.User, post.Forum)

	}
	if post.Parent ==  nil {
		if _, err := db.Map.Exec("UPDATE  post SET first_path = ? WHERE id = ? ", id,  id ); err != nil{
			println("FILLING FIRST PATH, err:", err)
		}
	}else{
		var tempPost Post
		db.Map.SelectOne(&tempPost, "SELECT first_path,last_path FROM post WHERE id=?", post.Parent)
		parentFirstPath := tempPost.FirstPath
		parentLastPath := tempPost.LastPath
		id_for_str := id
		if parentLastPath == "" {
			db.Map.Exec("UPDATE post SET first_path=?, last_path=? WHERE id=?", parentFirstPath, id_for_str ,  id)

		}else {
			parentLastPath += "."
			i := id
			var i64 int64
			i64 = int64(i)
			id_str:=strconv.FormatInt(i64, 10)
			db.Map.Exec("UPDATE post SET first_path=?, last_path=? WHERE id=?", parentFirstPath,id_str, id)

		}

	}



	db.Map.Exec("UPDATE thread SET posts = posts + 1 WHERE id = ?", post.Thread)
	context.JSON(200, gin.H{"code": 0, "response": gin.H{"date": post.Date, "forum": post.Forum, "id": id, "isApproved": post.IsApproved, "isDeleted": post.IsDeleted, "isEdited": post.IsEdited, "isHighlighted": post.IsHighlighted, "isSpam": post.IsSpam, "message": post.Message, "parent": post.Parent, "thread": post.Thread, "user": post.User}})
}

func (db *Resource) postDetails(context *gin.Context) {
	a := (context.Query("post"))
	post, err := strconv.Atoi(a)
	if err != nil{
		//println("fuck up")
	}
	if response := db.postWithID(post); response != nil {

		for _ , entity := range context.Request.URL.Query()["related"] {
			if entity == "user" {
				response[entity] = db.userWithEmail(response[entity].(string))
			} else if entity == "thread" {
				response[entity] = db.threadWithID(response[entity].(int))
			} else if entity == "forum" {
				response[entity] = db.forumWithShortName(response[entity].(string))
			}
		}
		context.JSON(200, gin.H{"code": 0, "response": response})
	} else {
		//println("fuck up")

		context.JSON(200, gin.H{"code": 1, "response": "Post not found"})
	}
}

func (db *Resource) postList(context *gin.Context) {
	query := "SELECT * FROM post WHERE"
	if forum := context.Query("forum"); forum != "" {
		query += " forum = " + "\"" + forum + "\""
	} else {
		query += " thread = " + context.Query("thread")
	}
	if since := context.Query("since"); since != "" {
		query += " AND date >= " + "\"" + since + "\""
	}
	query += " ORDER BY date " + context.DefaultQuery("order", "desc")
	if limit := context.Query("limit"); limit != "" {
		query += " LIMIT " + limit
	}
	var posts []Post
	_,err := db.Map.Select(&posts, query)
	print(err)
	context.JSON(200, gin.H{"code": 0, "response": posts})
}

func (db *Resource) postRemove(context *gin.Context) {
	var params struct {
		Post int `json:"post"`
	}
	context.BindJSON(&params)
	db.Map.Exec("UPDATE post SET isDeleted = true WHERE id = ?", params.Post)
	thread, _ := db.Map.SelectInt("SELECT thread FROM post WHERE id = ?", params.Post)
	db.Map.Exec("UPDATE thread SET posts = posts - 1 WHERE id = ?", thread)
	context.JSON(200, gin.H{"code": 0, "response": params})
}

func (db *Resource) postRestore(context *gin.Context) {
	var params struct {
		Post int `json:"post"`
	}
	context.BindJSON(&params)
	db.Map.Exec("UPDATE post SET isDeleted = false WHERE id = ?", params.Post)
	thread, _ := db.Map.SelectInt("SELECT thread FROM post WHERE id = ?", params.Post)
	db.Map.Exec("UPDATE thread SET posts = posts + 1 WHERE id = ?", thread)
	context.JSON(200, gin.H{"code": 0, "response": params})
}

func (db *Resource) postUpdate(context *gin.Context) {
	var params struct {
		Post    int    `json:"post"`
		Message string `json:"message"`
	}
	context.BindJSON(&params)
	db.Map.Exec("UPDATE post SET message = ? WHERE id = ?", params.Message, params.Post)
	context.JSON(200, gin.H{"code": 0, "response": db.postWithID(params.Post)})
}

func (db *Resource) postVote(context *gin.Context) {
	var params struct {
		Vote int `json:"vote"`
		Post int `json:"post"`
	}
	context.BindJSON(&params)
	if params.Vote > 0 {
		db.Map.Exec("UPDATE post SET likes = likes + 1, points = points + 1 WHERE id = ?", params.Post)
	} else {
		db.Map.Exec("UPDATE post SET dislikes = dislikes + 1, points = points - 1 WHERE id = ?", params.Post)
	}
	context.JSON(200, gin.H{"code": 0, "response": db.postWithID(params.Post)})
}
