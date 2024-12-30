package noti

import (
	"strconv"

	"github.com/21CLC01-WNC-Banking/WNC-Banking-BE/internal/domain/model"
)

func GenerateContent(req model.NotificationResponse) string {
	return `
		{
			"Name": ` + req.Name + `,
			"Amount": ` + strconv.Itoa(req.Amount) + `,
			"Transaction ID": ` + req.TransactionId + `,
			"Type": ` + req.Type + `,
			"Created At": ` + req.CreatedAt.String() + `
		}`
}
