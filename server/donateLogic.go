package server

import (
	"ecpayHook/redis"
	"ecpayHook/stringVerify"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"os"
	"path/filepath"
	"time"
	"unicode/utf8"
)

func donateLogic(c *gin.Context) {
	var dr donateRequest
	if err := c.BindJSON(&dr); err != nil {
		log.Println(err)
		c.JSON(400, gin.H{"status": false, "msg": "輸入的資料有誤"})
		return
	}

	if dr.TotalAmount < 65 || utf8.RuneCountInString(dr.Name) > 10 || utf8.RuneCountInString(dr.Name) == 0 || utf8.RuneCountInString(dr.Message) > 30 {
		c.JSON(400, gin.H{"status": false, "msg": "輸入的資料有誤"})
		return
	}

	if !stringVerify.StringVerify(dr.Name, dr.Message) {
		c.JSON(400, gin.H{"status": false, "msg": "輸入的資料中包含了敏感詞彙"})

		go func() {
			fileName := filepath.Join("log", fmt.Sprintf("%d.json", time.Now().Unix()))
			binary, err := dr.MarshalBinary()
			if err != nil {
				log.Println(err)
				return
			}
			err = os.WriteFile(fileName, binary, 0622)
			if err != nil {
				log.Println(err)
				return
			}
		}()

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
