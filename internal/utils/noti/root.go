package noti

import (
	"strconv"

	"github.com/21CLC01-WNC-Banking/WNC-Banking-BE/internal/domain/model"
)

func GenerateContentForTransactionNoti(req model.TransactionNotificationContent) string {
	return `
		{
			"Name": ` + req.Name + `,
			"Amount": ` + strconv.Itoa(req.Amount) + `,
			"Transaction ID": ` + req.TransactionId + `,
			"Type": ` + req.Type + `,
			"Created At": ` + req.CreatedAt.String() + `
		}`
}
