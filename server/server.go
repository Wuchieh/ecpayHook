package server

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"os"
	"strings"
)

var (
	setting Setting
)

type Setting struct {
	ServerAddr string `json:"serverAddr"`
	Mode       string `json:"mode"`

	MerchantID int    `json:"MerchantID"`
	HashKey    string `json:"HashKey"`
	HashIV     string `json:"HashIV"`
	DomainName string `json:"domainName"`
	LiveDebug  bool   `json:"liveDebug"`
}

func init() {
	if file, err := os.ReadFile("setting.json"); err != nil {
		panic(err)
	} else {
		err = json.Unmarshal(file, &setting)
		if err != nil {
			panic(err)
		}
	}
}

func Run() error {
	setMode()
	r := gin.Default()
	r.HTMLRender = createMyRender()
	r.Static("/api/static", "static")
	router(r)
	return r.Run(setting.ServerAddr)
}

func setMode() {
	switch strings.ToLower(setting.Mode[:1]) {
	case "r":
		gin.SetMode(gin.ReleaseMode)
	case "t":
		gin.SetMode(gin.TestMode)
	default:
		gin.SetMode(gin.DebugMode)
	}
}
