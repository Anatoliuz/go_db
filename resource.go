package main

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"gopkg.in/gorp.v1"
)

// Resource is a wrapper for db
type Resource struct {
	Map *gorp.DbMap
}

func openСonnection() *Resource {
	db, _ := sql.Open("mysql", "root:1@/TPForum?charset=utf8")
	dbmap := &gorp.DbMap{Db: db, Dialect: gorp.MySQLDialect{Engine: "InnoDB", Encoding: "UTF8"}}
	return &Resource{Map: dbmap}
}
