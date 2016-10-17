package main

import "github.com/gin-gonic/gin"

// Forum is a model
type Forum struct {
	ID        int    `db:"id" json:"id"`
	Name      string `db:"name" json:"name"`
	ShortName string `db:"short_name" json:"short_name"`
	User      string `db:"user" json:"user"`
}

func (db *Resource) forumWithShortName(shortName string) gin.H {
	var forum Forum
	db.Map.SelectOne(&forum, "SELECT * FROM forum WHERE short_name = ?", shortName)
	return gin.H{"id": forum.ID, "name": forum.Name, "short_name": forum.ShortName, "user": forum.User}
}

func (db *Resource) forumCreate(context *gin.Context) {
	var forum Forum
	context.BindJSON(&forum)
	db.Map.Exec("INSERT INTO forum (name, short_name, user) VALUES (?, ?, ?)", forum.Name, forum.ShortName, forum.User)
	db.Map.SelectOne(&forum, "SELECT * FROM forum WHERE short_name = ?", forum.ShortName)
	context.JSON(200, gin.H{"code": 0, "response": gin.H{"id": forum.ID, "name": forum.Name, "short_name": forum.ShortName, "user": forum.User}})
}

func (db *Resource) forumDetails(context *gin.Context) {
	response := db.forumWithShortName(context.Query("forum"))
	for _, entity := range context.Request.URL.Query()["related"] {
		if entity == "user" {
			response[entity] = db.userWithEmail(response[entity].(string))
		}
	}
	context.JSON(200, gin.H{"code": 0, "response": response})
}

func (db *Resource) forumListPosts(context *gin.Context) {
	related, relatedUser, relatedForum, relatedThread := context.Request.URL.Query()["related"], false, false, false
	for _, entity := range related {
		if entity == "user" {
			relatedUser = true
		} else if entity == "forum" {
			relatedForum = true
		} else if entity == "thread" {
			relatedThread = true
		}
	}
	query := "SELECT * FROM post WHERE forum = " + "\"" + context.Query("forum") + "\""
	if since := context.Query("since"); since != "" {
		query += " AND date >= " + "\"" + since + "\""
	}
	query += " ORDER BY date " + context.DefaultQuery("order", "DESC")
	if limit := context.Query("limit"); limit != "" {
		query += " LIMIT " + limit
	}
	var posts []Post
	db.Map.Select(&posts, query)
	response := make([]gin.H, len(posts))
	for index, post := range posts {
		response[index] = gin.H{"date": post.Date, "dislikes": post.Dislikes, "forum": post.Forum, "id": post.ID, "isApproved": post.IsApproved, "isDeleted": post.IsDeleted, "isEdited": post.IsEdited, "isHighlighted": post.IsHighlighted, "isSpam": post.IsSpam, "likes": post.Likes, "message": post.Message, "parent": post.Parent, "points": post.Points, "thread": post.Thread, "user": post.User}
		if relatedUser {
			response[index]["user"] = db.userWithEmail(response[index]["user"].(string))
		}
		if relatedForum {
			response[index]["forum"] = db.forumWithShortName(response[index]["forum"].(string))
		}
		if relatedThread {
			response[index]["thread"] = db.threadWithID(response[index]["thread"].(int))
		}
	}
	context.JSON(200, gin.H{"code": 0, "response": response})
}

func (db *Resource) forumListThreads(context *gin.Context) {
	related, relatedUser, relatedForum := context.Request.URL.Query()["related"], false, false
	for _, entity := range related {
		if entity == "user" {
			relatedUser = true
		} else if entity == "forum" {
			relatedForum = true
		}
	}
	query := "SELECT * FROM thread WHERE forum = " + "\"" + context.Query("forum") + "\""
	if since := context.Query("since"); since != "" {
		query += " AND date >= " + "\"" + since + "\""
	}
	query += " ORDER BY date " + context.DefaultQuery("order", "DESC")
	if limit := context.Query("limit"); limit != "" {
		query += " LIMIT " + limit
	}
	var threads []Thread
	db.Map.Select(&threads, query)
	response := make([]gin.H, len(threads))
	for index, thread := range threads {
		response[index] = gin.H{"date": thread.Date, "dislikes": thread.Dislikes, "forum": thread.Forum, "id": thread.ID, "isClosed": thread.IsClosed, "isDeleted": thread.IsDeleted, "likes": thread.Likes, "message": thread.Message, "points": thread.Points, "posts": thread.Posts, "slug": thread.Slug, "title": thread.Title, "user": thread.User}
		if relatedUser {
			response[index]["user"] = db.userWithEmail(response[index]["user"].(string))
		}
		if relatedForum {
			response[index]["forum"] = db.forumWithShortName(response[index]["forum"].(string))
		}
	}
	context.JSON(200, gin.H{"code": 0, "response": response})
}

func (db *Resource) forumListUsers(context *gin.Context) {
	query := "SELECT * FROM user WHERE email IN (SELECT DISTINCT user FROM post WHERE forum = " + "\"" + context.Query("forum") + "\")"
	if sinceID := context.Query("since_id"); sinceID != "" {
		query += " AND user.id >= " + sinceID
	}
	query += " ORDER BY user.name " + context.DefaultQuery("order", "DESC")
	if limit := context.Query("limit"); limit != "" {
		query += " LIMIT " + limit
	}
	var users []User
	db.Map.Select(&users, query)
	response := make([]gin.H, len(users))
	for index, user := range users {
		var followers []string
		db.Map.Select(&followers, "SELECT follower FROM follow WHERE followed = ?", user.Email)
		var following []string
		db.Map.Select(&following, "SELECT followed FROM follow WHERE follower = ?", user.Email)
		var subscriptions []int
		db.Map.Select(&subscriptions, "SELECT thread FROM subscription WHERE user = ?", user.Email)
		response[index] = gin.H{"about": user.About, "email": user.Email, "followers": followers, "following": following, "id": user.ID, "isAnonymous": user.IsAnonymous, "name": user.Name, "subscriptions": subscriptions, "username": user.Username}
	}
	context.JSON(200, gin.H{"code": 0, "response": response})
}
