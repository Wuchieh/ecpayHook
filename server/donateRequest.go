package server

import (
	"ecpayHook/database"
	"encoding/json"
	"fmt"
	"math/rand"
	"strconv"
	"time"
)

type donateRequest struct {
	Name                 string               `json:"Name"`
	Message              string               `json:"Message"`
	TotalAmount          int                  `json:"TotalAmount"`
	DonateTo             string               `json:"DonateTo"`
	PaymentStatus        bool                 `json:"PaymentStatus"`
	ProductionOrdersData ProductionOrdersData `json:"ProductionOrdersData"`
}

func (d *donateRequest) MarshalBinary() ([]byte, error) {
	return json.Marshal(d)
}

func (d *donateRequest) ToEcpayPaidData(simulatePaid bool) *database.EcpayPaidData {
	return database.NewEcpayPaidData(d.ProductionOrdersData.MerchantTradeNo, d.Name, d.TotalAmount, d.DonateTo, d.Message, simulatePaid)
}

func genNewOrdersData(TotalAmount int, TradeDesc string) ProductionOrdersData {
	now := time.Now()
	sprint := fmt.Sprint(now.UnixNano())
	MerchantTradeNo := fmt.Sprint(sprint + strconv.Itoa(rand.Intn(10)))
	var pod = ProductionOrdersData{
		MerchantTradeNo:   MerchantTradeNo,
		MerchantTradeDate: now.Format("2006/01/02 15:04:05"),
		TotalAmount:       TotalAmount,
		TradeDesc:         TradeDesc,
		ItemName:          "donate",
		ReturnURL:         getReturnURL(MerchantTradeNo),
		ChoosePayment:     "ALL",
		OrderResultURL:    getOrderResultURL(MerchantTradeNo),
		MerchantID:        setting.MerchantID,
		InvoiceMark:       "N",
		IgnorePayment:     "BARCODE",
		EncryptType:       "1",
		PaymentType:       "aio",
	}
	pod.CheckMacValue = pod.generateCheckMacValue()
	return pod
}

func getRedirectURL(id string) string {
	return fmt.Sprintf("https://%s/api/donate/%s", setting.DomainName, id)
}
