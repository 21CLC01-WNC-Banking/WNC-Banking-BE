package noti

import (
	"encoding/json"
	"fmt"

	"github.com/21CLC01-WNC-Banking/WNC-Banking-BE/internal/domain/model"
)

func GenerateContentForTransactionNoti(req model.TransactionNotificationContent) []byte {
	data := map[string]interface{}{
		"Name":          req.Name,
		"Amount":        req.Amount,
		"TransactionID": req.TransactionId,
		"Type":          req.Type,
		"CreatedAt":     req.CreatedAt,
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		fmt.Println(err)
	}

	return jsonData
}
