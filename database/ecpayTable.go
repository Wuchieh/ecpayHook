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

func NewEcpayPaidData(merchantTradeNo string, name string, totalAmount int, donateTo string, message string, simulatePaid bool) *EcpayPaidData {
	return &EcpayPaidData{MerchantTradeNo: merchantTradeNo, Name: name, TotalAmount: totalAmount, DonateTo: donateTo, Message: message, SimulatePaid: simulatePaid}
}

func GetEcpayPaidData(limit, page int) ([]EcpayPaidData, error) {
	query, err := db.Query("SELECT (id, merchant_trade_no, name, total_amount, donate_to, payment_time, merchant_trade_no, simulate_paid) FROM \"ecpay_paid_data\" LIMIT $1 OFFSET $2", limit, limit*(page-1))
	if err != nil {
		return nil, err
	}

	var ecpayPaidDatas []EcpayPaidData
	for query.Next() {
		var epd EcpayPaidData
		if err = query.Scan(
			&epd.Id,
			&epd.MerchantTradeNo,
			&epd.Name,
			&epd.TotalAmount,
			&epd.DonateTo,
			&epd.PaymentTime,
			&epd.Message,
			&epd.SimulatePaid); err != nil {
			continue
		}
		ecpayPaidDatas = append(ecpayPaidDatas, epd)
	}

	return ecpayPaidDatas, err
}

func (e *EcpayPaidData) Save() error {
	tx, err := db.Begin()
	if err != nil {
		return pgError(err)
	}
	defer tx.Rollback()

	stmt, err := db.Prepare("INSERT INTO \"ecpay_paid_data\" (merchant_trade_no, name, total_amount, donate_to, message, simulate_paid) VALUES ($1,$2,$3,$4,$5,$6)")
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
