package v1

import (
	"github.com/21CLC01-WNC-Banking/WNC-Banking-BE/internal/domain/entity"
	httpcommon "github.com/21CLC01-WNC-Banking/WNC-Banking-BE/internal/domain/http_common"
	"github.com/21CLC01-WNC-Banking/WNC-Banking-BE/internal/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

type DebtReplyHandler struct {
	debtReplyService service.DebtReplyService
}

func NewDebtReplyHandler(debtReplyService service.DebtReplyService) *DebtReplyHandler {
	return &DebtReplyHandler{
		debtReplyService: debtReplyService,
	}
}

// @Summary User get debt reply
// @Description User get debt reply by a debt reminder id
// @Tags Debt Reminder
// @Produce  json
// @Param debtReminderId path string true "debtReminder Id"
// @Router /debt-reply/{debtReminderId} [get]
// @Success 200 {object} httpcommon.HttpResponse[entity.DebtReply]
// @Failure 400 {object} httpcommon.HttpResponse[any]
// @Failure 500 {object} httpcommon.HttpResponse[any]
func (handler *DebtReplyHandler) GetReplyByDebtReminderId(ctx *gin.Context) {
	debtReminderId := ctx.Param("debtReminderId")
	debtReply, err := handler.debtReplyService.GetReplyByDebtReminderId(ctx, debtReminderId)
	if err != nil {
		if err.Error() == httpcommon.ErrorMessage.SqlxNoRow {
			ctx.JSON(http.StatusNotFound, httpcommon.NewErrorResponse(
				httpcommon.Error{Message: err.Error(), Code: httpcommon.ErrorResponseCode.RecordNotFound},
			))
			return
		}
		ctx.JSON(http.StatusInternalServerError, httpcommon.NewErrorResponse(
			httpcommon.Error{Message: err.Error(), Code: httpcommon.ErrorResponseCode.InternalServerError},
		))
		return
	}
	ctx.JSON(http.StatusOK, httpcommon.NewSuccessResponse[entity.DebtReply](debtReply))
}
