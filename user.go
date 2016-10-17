package main

import "github.com/gin-gonic/gin"

// User is a model
type User struct {
	About       *string `db:"about" json:"about"`
	Email       string  `db:"email" json:"email"`
	ID          int     `db:"id" json:"id"`
	IsAnonymous bool    `db:"isAnonymous" json:"isAnonymous"`
	Name        *string `db:"name" json:"name"`
	Username    *string `db:"username" json:"username"`
}

func (db *Resource) userWithEmail(email string) gin.H {
	var user User
	db.Map.SelectOne(&user, "SELECT * FROM user WHERE email = ?", email)
	var followers []string
	db.Map.Select(&followers, "SELECT follower FROM follow WHERE followed = ?", email)
	var following []string
	db.Map.Select(&following, "SELECT followed FROM follow WHERE follower = ?", email)
	subscriptions := make([]int, 0)
	db.Map.Select(&subscriptions, "SELECT thread FROM subscription WHERE user = ?", email)
	return gin.H{"about": user.About, "email": user.Email, "followers": followers, "following": following, "id": user.ID, "isAnonymous": user.IsAnonymous, "name": user.Name, "subscriptions":subscriptions,  "username": user.Username}
}

func (db *Resource) userCreate(context *gin.Context) {
	var user User
	context.BindJSON(&user)
	if result, err := db.Map.Exec("INSERT INTO user (about, email, isAnonymous, name, username) VALUES (?, ?, ?, ?, ?)", user.About, user.Email, user.IsAnonymous, user.Name, user.Username); err == nil {
		id, _ := result.LastInsertId()
		context.JSON(200, gin.H{"code": 0, "response": gin.H{"about": user.About, "email": user.Email, "id": id, "isAnonymous": user.IsAnonymous, "name": user.Name, "username": user.Username}})
	} else {
		context.JSON(200, gin.H{"code": 5, "response": "User already exists"})
	}
}

func (db *Resource) userDetails(context *gin.Context) {
	context.JSON(200, gin.H{"code": 0, "response": db.userWithEmail(context.Query("user"))})
}

func (db *Resource) userFollow(context *gin.Context) {
	var params struct {
		Follower string `json:"follower"`
		Followee string `json:"followee"`
	}
	context.BindJSON(&params)
	db.Map.Exec("INSERT INTO follow (follower, followee) VALUES (?, ?)", params.Follower, params.Followee)
	context.JSON(200, gin.H{"code": 0, "response": db.userWithEmail(params.Follower)})
}

func (db *Resource) userListFollowers(context *gin.Context) {
	query := "SELECT follower FROM follow JOIN user ON follow.follower = user.email WHERE followee = " + "\"" + context.Query("user") + "\""
	if sinceID := context.Query("since_id"); sinceID != "" {
		query += " AND id >= " + sinceID
	}
	query += " ORDER BY follower " + context.DefaultQuery("order", "desc")
	if limit := context.Query("limit"); limit != "" {
		query += " LIMIT " + limit
	}
	var emails []string
	db.Map.Select(&emails, query)
	response := make([]gin.H, len(emails))
	for index, email := range emails {
		response[index] = db.userWithEmail(email)
	}
	context.JSON(200, gin.H{"code": 0, "response": response})
}

func (db *Resource) userListFollowing(context *gin.Context) {
	query := "SELECT followed FROM follow JOIN user ON follow.followed = user.email WHERE follower = " + "\"" + context.Query("user") + "\""
	if sinceID := context.Query("since_id"); sinceID != "" {
		query += " AND id >= " + sinceID
	}
	query += " ORDER BY followed " + context.DefaultQuery("order", "desc")
	if limit := context.Query("limit"); limit != "" {
		query += " LIMIT " + limit
	}
	var emails []string
	db.Map.Select(&emails, query)
	response := make([]gin.H, len(emails))
	for index, email := range emails {
		response[index] = db.userWithEmail(email)
	}
	context.JSON(200, gin.H{"code": 0, "response": response})
}

func (db *Resource) userListPosts(context *gin.Context) {
	query := "SELECT * FROM post WHERE user = " + "\"" + context.Query("user") + "\""
	if since := context.Query("since"); since != "" {
		query += " AND date >= " + "\"" + since + "\""
	}
	query += " ORDER BY date " + context.DefaultQuery("order", "desc")
	if limit := context.Query("limit"); limit != "" {
		query += " LIMIT " + limit
	}
	var posts []Post
	db.Map.Select(&posts, query)
	context.JSON(200, gin.H{"code": 0, "response": posts})
}

func (db *Resource) userUnfollow(context *gin.Context) {
	var params struct {
		Follower string `json:"follower"`
		Followee string `json:"followed"`
	}
	context.BindJSON(&params)
	db.Map.Exec("DELETE FROM follow WHERE follower = ? AND followed = ?", params.Follower, params.Followee)
	context.JSON(200, gin.H{"code": 0, "response": db.userWithEmail(params.Follower)})
}

func (db *Resource) userUpdateProfile(context *gin.Context) {
	var params struct {
		About string `json:"about"`
		User  string `json:"user"`
		Name  string `json:"name"`
	}
	context.BindJSON(&params)
	db.Map.Exec("UPDATE user SET about = ?, name = ? WHERE email = ?", params.About, params.Name, params.User)
	context.JSON(200, gin.H{"code": 0, "response": db.userWithEmail(params.User)})
}
