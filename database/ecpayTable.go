package database

import (
	"github.com/google/uuid"
	"time"
)

type EcpayPaidData struct {
	Id              uuid.UUID
	MerchantTradeNo string
	Name            string
	TotalAmount     int
	DonateTo        string
	PaymentTime     time.Time
	Message         string
	SimulatePaid    bool
}

func NewEcpayTable(merchantTradeNo string, name string, totalAmount int, donateTo string, message string, simulatePaid bool) *EcpayPaidData {
	return &EcpayPaidData{MerchantTradeNo: merchantTradeNo, Name: name, TotalAmount: totalAmount, DonateTo: donateTo, Message: message, SimulatePaid: simulatePaid}
}

func (e *EcpayPaidData) Save() error {
	tx, err := db.Begin()
	if err != nil {
		return pgError(err)
	}
	defer tx.Rollback()

	stmt, err := db.Prepare("INSERT INTO ecpay_paid_data (merchant_trade_no, name, total_amount, donate_to, message, simulate_paid) VALUES ($1,$2,$3,$4,$5,$6)")
	if err != nil {
		return pgError(err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(e.MerchantTradeNo,
		e.Name,
		e.TotalAmount,
		e.DonateTo,
		e.Message,
		e.SimulatePaid)
	if err != nil {
		return pgError(err)
	}

	err = tx.Commit()
	if err != nil {
		return pgError(err)
	}

	return nil
}
