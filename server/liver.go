package server

import (
	"ecpayHook/redis"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"log"
	"os"
)

var (
	livers         = make(map[string]string)
	wsConnectUsers = make(map[string][]*websocket.Conn)
)

func init() {
	if file, err := os.ReadFile("users.json"); err != nil {
		panic(err)
	} else {
		if err = json.Unmarshal(file, &livers); err != nil {
			panic(err)
		}
	}
}

// 發送通知
func sendAlert(id string) {
	go func() {
		get := redis.Get("donate_" + id)
		if err := get.Err(); err != nil {
			log.Println(err)
			return
		}
		redisByte, err := get.Bytes()
		if err != nil {
			log.Println(err)
			return
		}
		var dr donateRequest
		err = json.Unmarshal(redisByte, &dr)
		if err != nil {
			log.Println(err)
			return
		}

		if _, ok := wsConnectUsers[dr.DonateTo]; ok {
			for _, conn := range wsConnectUsers[dr.DonateTo] {
				_ = conn.WriteJSON(dr)
			}
		}

		if err != nil {
			log.Println(err)
			return
		}
	}()
}

func showLiverAlert(c *gin.Context) {
	password := c.Param("password")[1:]
	fmt.Println(password)

	var user string
	for k, v := range livers {
		if v == password {
			user = k
		}
	}
	fmt.Println(user)

	c.HTML(200, "showLiverAlert", gin.H{"user": user})
}

func wsLiverAlert(c *gin.Context) {
	user := c.Param("user")

	// 建立連線
	ws, err := genWebSocket(c)
	if err != nil {
		log.Println(err)
		return
	}

	// 連線關閉
	defer func(ws *websocket.Conn) {
		wsClose(ws)
		for i, conn := range wsConnectUsers[user] {
			if conn == ws {
				wsConnectUsers[user] = append(wsConnectUsers[user][:i], wsConnectUsers[user][i+1:]...)
				return
			}
		}
	}(ws)

	if _, ok := wsConnectUsers[user]; ok {
		wsConnectUsers[user] = append(wsConnectUsers[user], ws)
	} else {
		wsConnectUsers[user] = []*websocket.Conn{}
		wsConnectUsers[user] = append(wsConnectUsers[user], ws)
	}

	//_ = ws.WriteMessage(websocket.TextMessage, []byte(user+" 你好"))

	// 開始工作
	wsProcess(ws)
}
