package main

import (
	"strconv"
	"github.com/gin-gonic/gin"
	"sort"
)

// Thread is a model
type Thread struct {
	Date      string `db:"date" json:"date"`
	Dislikes  int    `db:"dislikes" json:"dislikes"`
	Forum     string `db:"forum" json:"forum"`
	ID        int    `db:"id" json:"id"`
	IsClosed  bool   `db:"isClosed" json:"isClosed"`
	IsDeleted bool   `db:"isDeleted" json:"isDeleted"`
	Likes     int    `db:"likes" json:"likes"`
	Message   string `db:"message" json:"message"`
	Points    int    `db:"points" json:"points"`
	Posts     int    `db:"posts" json:"posts"`
	Slug      string `db:"slug" json:"slug"`
	Title     string `db:"title" json:"title"`
	User      string `db:"user" json:"user"`
}

func (db *Resource) threadWithID(id int) gin.H {
	var thread Thread
	db.Map.SelectOne(&thread, "SELECT * FROM thread WHERE id = ?", id)
	return gin.H{"date": thread.Date, "dislikes": thread.Dislikes, "forum": thread.Forum, "id": thread.ID, "isClosed": thread.IsClosed, "isDeleted": thread.IsDeleted, "likes": thread.Likes, "message": thread.Message, "points": thread.Points, "posts": thread.Posts, "slug": thread.Slug, "title": thread.Title, "user": thread.User}
}
func (db *Resource) threadClose(context *gin.Context) {
	var params struct {
		Thread int `json:"thread"`
	}
	context.BindJSON(&params)
	db.Map.Exec("UPDATE thread SET isClosed = true WHERE id = ?", params.Thread)
	context.JSON(200, gin.H{"code": 0, "response": params})
}

func (db *Resource) threadCreate(context *gin.Context) {
	var thread Thread
	context.BindJSON(&thread)
	result, err := db.Map.Exec("INSERT INTO thread (date, forum, isClosed, isDeleted, message, slug, title, user) VALUES (?, ?, ?, ?, ?, ?, ?, ?)", thread.Date, thread.Forum, thread.IsClosed, thread.IsDeleted, thread.Message, thread.Slug, thread.Title, thread.User)
	println("err:", err)
	id, _ := result.LastInsertId()
	context.JSON(200, gin.H{"code": 0, "response": gin.H{"date": thread.Date, "forum": thread.Forum, "id": id, "isClosed": thread.IsClosed, "isDeleted": thread.IsDeleted, "message": thread.Message, "slug": thread.Slug, "title": thread.Title, "user": thread.User}})
}

func (db *Resource) threadDetails(context *gin.Context) {
	thread, _ := strconv.Atoi(context.Query("thread"))
	response := db.threadWithID(thread)
	for _, entity := range context.Request.URL.Query()["related"] {
		if entity == "user" {
			response[entity] = db.userWithEmail(response[entity].(string))
		} else if entity == "forum" {
			response[entity] = db.forumWithShortName(response[entity].(string))
		} else {
			context.JSON(200, gin.H{"code": 3, "response": "Bad request"})
			return
		}
	}
	context.JSON(200, gin.H{"code": 0, "response": response})
}

func (db *Resource) threadList(context *gin.Context) {
	query := "SELECT * FROM thread WHERE"
	if user := context.Query("user"); user != "" {
		query += " user = " + "\"" + user + "\""
		} else {
		query += " forum  = " + "\"" + context.Query("forum") + "\""
	}
	if since := context.Query("since"); since != "" {
		query += " AND date >= " + "\"" + since + "\""
	}
	query += " ORDER BY date " + context.DefaultQuery("order", "desc")
	if limit := context.Query("limit"); limit != "" {
		query += " LIMIT " + limit
	}
	var threads []Thread
	db.Map.Select(&threads, query)
	context.JSON(200, gin.H{"code": 0, "response": threads})
}

func (db *Resource) threadListPosts(context *gin.Context) {
	var posts []Post

		query := "SELECT * FROM post WHERE thread = " + context.Query("thread")
		if since := context.Query("since"); since != "" {
			query += " AND date >= " + "\"" + since + "\""
		}

		sortType := context.Query("sort")
	        if sortType == "" {
			query += " ORDER BY date " + context.DefaultQuery("order", "desc")
			if limit := context.Query("limit"); limit != "" {
				query += " LIMIT " + limit
			}
		}else if sortType == "flat" {
			query += " ORDER BY date " + context.DefaultQuery("order", "desc")
			if limit := context.Query("limit"); limit != "" {
				query += " LIMIT " + limit
			}
		}else {

			//query := "SELECT * FROM post WHERE thread = " + context.Query("thread")
			//if since := context.Query("since"); since != "" {
			//	query += " AND date >= " + "\"" + since + "\""
			//}
			//query += " ORDER BY date " + context.DefaultQuery("order", "desc")
			if limit := context.Query("limit"); limit != "" {
				query += " LIMIT " + limit
			}

			//if order := context.Query("order"); order == "asc" {
			//	sort.Sort(FirstPathASC(FirstPathASC(posts)))
			//	sort.Sort(LastPathASC(LastPathASC(posts)))
			//}
			//order := context.Query("order")



		}
		//if err != nil{
		//	println(err)
		//}
		//println("bef")
		//for _, post := range posts {
		//	println(post.Forum)
		//}
		//println("after")
		//
		//sort.Sort(ByAge(posts))
		//for _, post := range posts {
		//	println(post.Forum)
		//}
		//if sort := context.Query("sort"); sort == "tree" {
		//	query += " ORDER BY first_path " + context.DefaultQuery("order", "desc")
		//	if limit := context.Query("limit"); limit != "" {
		//		query += " LIMIT " + limit
		//	}
		//}
		//if sort := context.Query("sort"); sort == "parent_tree"{
		//	if limit := context.Query("limit"); limit != "" {
		//		//query += " LIMIT " +
		//	err,_ := db.Map.Exec("SELECT p1.* FROM post AS p1 WHERE p1.date >= " + "\"" +  context.Query("since") + "\"" + " AND first_path IN ( SELECT * FROM ( SELECT DISTINCT first_path FROM post WHERE post.id = %d ORDER BY first_index DESC LIMIT limit  ) AS p2) ORDER BY  " + "\"" +  context.Query("limit") + "\"" + " DESC , last_path")
		//		println(err)
		//	}
		//	db.Map.
		//}
		//db.Map.Select(&posts, query)
		//	context.JSON(200, gin.H{"code": 0, "response": posts})

	db.Map.Select(&posts, query)
	for _, p := range posts{
		print(p.FirstPath)
		println(p.LastPath)
	}

	sort.Sort(FirstPathDESC(FirstPathDESC(posts)))
	sort.Sort(LastPathDESC(LastPathDESC(posts)))
	context.JSON(200, gin.H{"code": 0, "response": posts})
}

func (db *Resource) threadOpen(context *gin.Context) {
	var params struct {
		Thread int `json:"thread"`
	}
	context.BindJSON(&params)
	db.Map.Exec("UPDATE thread SET isClosed = false WHERE id = ?", params.Thread)
	context.JSON(200, gin.H{"code": 0, "response": params})
}

func (db *Resource) threadRemove(context *gin.Context) {
	var params struct {
		Thread int `json:"thread"`
	}
	context.BindJSON(&params)
	db.Map.Exec("UPDATE thread SET isDeleted = true, posts = 0 WHERE id = ?", params.Thread)
	db.Map.Exec("UPDATE post SET isDeleted = true WHERE thread = ?", params.Thread)
	context.JSON(200, gin.H{"code": 0, "response": params})
}

func (db *Resource) threadRestore(context *gin.Context) {
	var params struct {
		Thread int `json:"thread"`
	}
	context.BindJSON(&params)
	posts, _ := db.Map.SelectInt("SELECT count(id) FROM post WHERE thread = ?", params.Thread)
	db.Map.Exec("UPDATE thread SET isDeleted = false, posts = ? WHERE id = ?", posts, params.Thread)
	db.Map.Exec("UPDATE post SET isDeleted = false WHERE thread = ?", params.Thread)
	context.JSON(200, gin.H{"code": 0, "response": params})
}

func (db *Resource) threadSubscribe(context *gin.Context) {
	var params struct {
		User   string `json:"user"`
		Thread int    `json:"thread"`
	}
	context.BindJSON(&params)
	db.Map.Exec("INSERT INTO subscription (user, thread) VALUES (?, ?)", params.User, params.Thread)
	context.JSON(200, gin.H{"code": 0, "response": params})
}

func (db *Resource) threadUnsubscribe(context *gin.Context) {
	var params struct {
		User   string `json:"user"`
		Thread int    `json:"thread"`
	}
	context.BindJSON(&params)
	db.Map.Exec("DELETE FROM subscription WHERE user = ? AND thread = ?", params.User, params.Thread)
	context.JSON(200, gin.H{"code": 0, "response": params})
}

func (db *Resource) threadUpdate(context *gin.Context) {
	var params struct {
		Message string `json:"message"`
		Slug    string `json:"slug"`
		Thread  int    `json:"thread"`
	}
	context.BindJSON(&params)
	db.Map.Exec("UPDATE thread SET message = ?, slug = ? WHERE id = ?", params.Message, params.Slug, params.Thread)
	context.JSON(200, gin.H{"code": 0, "response": db.threadWithID(params.Thread)})
}

func (db *Resource) threadVote(context *gin.Context) {
	var params struct {
		Vote   int `json:"vote"`
		Thread int `json:"thread"`
	}
	context.BindJSON(&params)
	if params.Vote > 0 {
		db.Map.Exec("UPDATE thread SET likes = likes + 1, points = points + 1 WHERE id = ?", params.Thread)
	} else {
		db.Map.Exec("UPDATE thread SET dislikes = dislikes + 1, points = points - 1 WHERE id = ?", params.Thread)
	}
	context.JSON(200, gin.H{"code": 0, "response": db.threadWithID(params.Thread)})
}
