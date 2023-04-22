package database

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/lib/pq"
	"log"
	"os"
)

var (
	db *sql.DB

	host     string
	user     string
	password string
	dbname   string
)

type setting struct {
	DatabaseHost     string `json:"databaseHost"`
	DatabaseUser     string `json:"databaseUser"`
	DatabasePassword string `json:"databasePassword"`
	DatabaseDbname   string `json:"databaseDbname"`
}

func init() {
	var s setting

	if file, err := os.ReadFile("setting.json"); err != nil {
		log.Panicln("database os 找不到設定檔 Error", err)
	} else {
		err = json.Unmarshal(file, &s)
		if err != nil {
			log.Panicln("database json 解析Error", err)
		}
	}

	host = s.DatabaseHost
	user = s.DatabaseUser
	password = s.DatabasePassword
	dbname = s.DatabaseDbname
}

func DatabaseInit() (err error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable", host, user, password, dbname)
	if db, err = sql.Open("postgres", dsn); err != nil {
		log.Println("sql.Open Error", err)
		return
	}
	if err = db.Ping(); err != nil {
		log.Println("db.Ping Error", err)
		return
	}
	return
}
