package server

import (
	"ecpayHook/redis"
	"ecpayHook/stringVerify"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"time"
)

func donateLogic(c *gin.Context) {
	var dr donateRequest
	if err := c.BindJSON(&dr); err != nil || dr.TotalAmount < 65 || len(dr.Name) > 10 || len(dr.Name) == 0 || len(dr.Message) > 30 {
		log.Println(err)
		c.JSON(400, gin.H{"status": false, "msg": "輸入的資料有誤"})
		return
	}

	if !stringVerify.StringVerify(dr.Message) {
		c.JSON(400, gin.H{"status": false, "msg": "輸入的資料中包含了敏感詞彙"})
		return
	}

	OrdersData := genNewOrdersData(dr.TotalAmount, "交易敘述")
	dr.ProductionOrdersData = OrdersData

	binary, err := dr.MarshalBinary()
	if err != nil {
		c.JSON(400, gin.H{"status": false, "msg": err.Error()})
		return
	}

	set := redis.Set(fmt.Sprintf("donate_%s", dr.ProductionOrdersData.MerchantTradeNo), binary, time.Hour*24*7)
	if set.Err() != nil {
		log.Println(set.Err())
		c.JSON(400, gin.H{"status": false, "msg": set.Err().Error()})
		return
	}

	c.JSON(200, gin.H{"status": true, "redirectURL": getRedirectURL(OrdersData.MerchantTradeNo)})
}
