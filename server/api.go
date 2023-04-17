package server

import (
	"crypto/sha256"
	"ecpayHook/redis"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"log"
	"net/url"
	"sort"
	"strings"
	"time"
)

func returnOrder(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		return
	}
	if p, ok := getReturnOrderData(c); ok {
		OrderProcessingLogic(id, p)
		c.String(200, "1|OK")
	} else {
		log.Println("接收到偽造訊息")
		c.String(200, "接收到偽造訊息")
	}
}

func OrderProcessingLogic(id string, p map[string]string) {
	if SimulatePaid, ok := p["SimulatePaid"]; ok {
		switch SimulatePaid {
		case "0":
			fmt.Println("非模擬付款")
			changeOrderStatus(id)
			sendAlert(id)
		case "1":
			fmt.Println("模擬付款")
			if setting.LiveDebug {
				sendAlert(id)
			}
		}
	}
}

// 解析ReturnOrderData
func getReturnOrderData(c *gin.Context) (map[string]string, bool) {
	bodyBytes, err := io.ReadAll(c.Request.Body)
	if err != nil {
		log.Println(err)
		return nil, false
	}

	values, err := url.ParseQuery(string(bodyBytes))
	if err != nil {
		log.Println(err)
		return nil, false
	}

	resultMap := make(map[string]string)
	for key, values := range values {
		if len(values) > 0 {
			resultMap[key] = values[0]
		}
	}
	if verifyReturnOrderData(resultMap) {
		return resultMap, true
	} else {
		return nil, false
	}
}
func verifyReturnOrderData(r map[string]string) bool {
	values := url.Values{}

	var keys []string
	for k, v := range r {
		if k != "CheckMacValue" {
			keys = append(keys, k)
			values.Add(k, v)
		}
	}
	sort.Strings(keys)

	// 串連參數
	var paramStr string
	for _, k := range keys {
		paramStr += k + "=" + values.Get(k) + "&"
	}

	// 加上HashKey及HashIV
	paramStr = "HashKey=" + setting.HashKey + "&" + paramStr + "HashIV=" + setting.HashIV

	// URL encode
	encodedStr := url.QueryEscape(paramStr)

	// 轉為小寫
	lowerStr := strings.ToLower(encodedStr)

	// 產生雜湊值
	h := sha256.New()
	h.Write([]byte(lowerStr))
	hash := hex.EncodeToString(h.Sum(nil))

	return strings.ToUpper(hash) == r["CheckMacValue"]
}

// 變更訂單狀態
func changeOrderStatus(id string) {
	get := redis.Get(fmt.Sprintf("donate_%s", id))
	if get.Err() != nil {
		log.Println(get.Err())
		return
	}
	bytes, err := get.Bytes()
	if err != nil {
		log.Println(get.Err())
		return
	}
	var dr donateRequest
	if err := json.Unmarshal(bytes, &dr); err != nil {
		log.Println(err)
		return
	}
	dr.PaymentStatus = true
	binary, err := dr.MarshalBinary()
	if err != nil {
		log.Println(err)
		return
	}
	redis.Set(fmt.Sprintf("donate_%s", id), binary, 5*time.Minute)
	fmt.Println(fmt.Sprintf("感謝 %s 贊助 %s %d元\n%s", dr.Name, dr.DonateTo, dr.ProductionOrdersData.TotalAmount, dr.Message))
}

func ResultOrder(c *gin.Context) {
	id := c.Param("id")
	get := redis.Get(fmt.Sprintf("donate_%s", id))
	if get.Err() != nil {
		c.String(400, get.Err().Error())
		return
	}
	bytes, err := get.Bytes()
	if err != nil {
		c.String(400, err.Error())
		return
	}
	var dr donateRequest
	if err = json.Unmarshal(bytes, &dr); err != nil {
		c.String(400, err.Error())
		return
	}
	if dr.PaymentStatus {
		c.String(200, "付款成功")
	} else {
		c.String(200, "付款尚未完成")
	}
}
