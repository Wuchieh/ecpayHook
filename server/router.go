package server

import (
	"github.com/gin-gonic/gin"
)

func router(r *gin.Engine) {
	r.GET("/", index)
	r.POST("/donate", donate)
	r.GET("/donate/:id", RedirectDonate)

	api := r.Group("/api")
	api.POST("/return/:id", returnOrder)
	api.GET("/result/:id", ResultOrder)
	api.POST("/donate", donate)
	api.GET("/donate/:id", RedirectDonate)

	liver := api.Group("/liver")
	liver.GET("/showLiverAlert/*password", showLiverAlert)
	liver.GET("/ws/:user", wsLiverAlert)
}

func _(c *gin.Context) {
	c.String(200, "1|OK")
}
