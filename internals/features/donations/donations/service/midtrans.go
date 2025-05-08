package service

import (
	"quizku/internals/features/donations/donations/model"

	midtrans "github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/snap"
)

var SnapClient snap.Client

func InitMidtrans(serverKey string) {
	SnapClient.New(serverKey, midtrans.Sandbox)
}

func GenerateSnapToken(d model.Donation, name string, email string) (string, error) {
	req := &snap.Request{
		TransactionDetails: midtrans.TransactionDetails{
			OrderID:  d.OrderID,
			GrossAmt: int64(d.Amount),
		},
		CustomerDetail: &midtrans.CustomerDetails{
			FName: name,
			Email: email,
		},
	}
	resp, err := SnapClient.CreateTransaction(req)
	if err != nil {
		return "", err
	}
	return resp.Token, nil
}
