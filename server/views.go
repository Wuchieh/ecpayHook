package server

import (
	"ecpayHook/redis"
	"encoding/json"
	"fmt"
	"github.com/gin-contrib/multitemplate"
	"github.com/gin-gonic/gin"
)

func createMyRender() multitemplate.Renderer {
	r := multitemplate.NewRenderer()
	r.AddFromFiles("index", "templates/index.gohtml")
	r.AddFromFiles("donate", "templates/donate.gohtml")
	r.AddFromFiles("showLiverAlert", "templates/showLiverAlert.gohtml")
	for k := range livers {
		fmt.Println(k)
		r.AddFromFiles(k, fmt.Sprintf("templates/%s.gohtml", k))
	}
	return r
}

func index(c *gin.Context) {
	c.HTML(200, "index", gin.H{"users": livers})
}

func donate(c *gin.Context) {
	donateLogic(c)
}

func RedirectDonate(c *gin.Context) {
	id := c.Param("id")
	get := redis.Get(fmt.Sprintf("donate_%s", id))
	bytes, err := get.Bytes()
	if err != nil {
		c.JSON(400, gin.H{"status": false, "msg": err.Error()})
		return
	}

	var dr donateRequest

	if err = json.Unmarshal(bytes, &dr); err != nil {
		c.JSON(400, gin.H{"status": false, "msg": err.Error()})
		return
	}
	c.HTML(200, "donate", dr.ProductionOrdersData)
}
