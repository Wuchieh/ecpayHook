package server

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/url"
	"sort"
	"strconv"
	"strings"
)

type ProductionOrdersData struct {
	ChoosePayment     string `json:"ChoosePayment"`
	EncryptType       string `json:"EncryptType"`
	IgnorePayment     string `json:"IgnorePayment"`
	InvoiceMark       string `json:"InvoiceMark"`
	ItemName          string `json:"ItemName"`
	MerchantID        int    `json:"MerchantID"`
	MerchantTradeDate string `json:"MerchantTradeDate"`
	MerchantTradeNo   string `json:"MerchantTradeNo"`
	OrderResultURL    string `json:"OrderResultURL"`
	PaymentType       string `json:"PaymentType"`
	ReturnURL         string `json:"ReturnURL"`
	TotalAmount       int    `json:"TotalAmount"`
	TradeDesc         string `json:"TradeDesc"`

	CheckMacValue string `json:"CheckMacValue"`
}

func (data ProductionOrdersData) generateCheckMacValue() string {
	values := url.Values{}

	values.Add("ChoosePayment", data.ChoosePayment)
	values.Add("EncryptType", data.EncryptType)
	values.Add("IgnorePayment", data.IgnorePayment)
	values.Add("InvoiceMark", data.InvoiceMark)
	values.Add("ItemName", data.ItemName)
	values.Add("MerchantID", strconv.Itoa(data.MerchantID))
	values.Add("MerchantTradeDate", data.MerchantTradeDate)
	values.Add("MerchantTradeNo", data.MerchantTradeNo)
	values.Add("OrderResultURL", data.OrderResultURL)
	values.Add("PaymentType", data.PaymentType)
	values.Add("ReturnURL", data.ReturnURL)
	values.Add("TotalAmount", strconv.Itoa(data.TotalAmount))
	values.Add("TradeDesc", data.TradeDesc)

	var keys []string
	for k := range values {
		keys = append(keys, k)
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

	// 轉為大寫
	return strings.ToUpper(hash)
}

func getReturnURL(id string) string {
	return fmt.Sprintf("https://%s/api/return/%s", setting.DomainName, id)
}

func getOrderResultURL(id string) string {
	return fmt.Sprintf("https://%s/api/result/%s", setting.DomainName, id)
}
