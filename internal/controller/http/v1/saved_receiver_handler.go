package v1

import (
	"net/http"
	"strconv"

	httpcommon "github.com/21CLC01-WNC-Banking/WNC-Banking-BE/internal/domain/http_common"
	"github.com/21CLC01-WNC-Banking/WNC-Banking-BE/internal/domain/model"
	"github.com/21CLC01-WNC-Banking/WNC-Banking-BE/internal/service"
	"github.com/21CLC01-WNC-Banking/WNC-Banking-BE/internal/utils/validation"
	"github.com/gin-gonic/gin"
)

type SavedReceiverHandler struct {
	savedReceiverService service.SavedReceiverService
}

func NewSavedReceiverHandler(savedReceiverService service.SavedReceiverService) *SavedReceiverHandler {
	return &SavedReceiverHandler{
		savedReceiverService: savedReceiverService,
	}
}

// @Summary Add Receiver
// @Description Add a new receiver
// @Tags Receivers
// @Accept  json
// @Produce  json
// @Param receiver body model.Receiver true "Receiver Payload"
// @Success 204 "No Content"
// @Failure 500 {object} httpcommon.HttpResponse[any] "Internal Server Error"
// @Router /customer/saved-receiver [post]
func (handler *SavedReceiverHandler) AddReceiver(ctx *gin.Context) {
	var receiver model.Receiver

	if err := validation.BindJsonAndValidate(ctx, &receiver); err != nil {
		return
	}
	err := handler.savedReceiverService.AddReceiver(ctx, receiver)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, httpcommon.NewErrorResponse(httpcommon.Error{
			Message: err.Error(), Field: "", Code: httpcommon.ErrorResponseCode.InternalServerError,
		}))
		return
	}
	ctx.AbortWithStatus(204)
}

// GetAllReceivers handles the endpoint to get all saved receivers
// @Summary Get all saved receivers
// @Description Fetches all saved receivers for the authenticated user
// @Tags Receivers
// @Accept  json
// @Produce  json
// @Success 200 {object} httpcommon.HttpResponse[[]model.SavedReceiverResponse]
// @Failure 500 {object} httpcommon.HttpResponse[any] "Internal Server Error"
// @Router /customer/saved-receiver [get]
func (handler *SavedReceiverHandler) GetAllReceivers(ctx *gin.Context) {
	savedReceivers, err := handler.savedReceiverService.GetAllReceivers(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, httpcommon.NewErrorResponse(httpcommon.Error{
			Message: err.Error(), Field: "", Code: httpcommon.ErrorResponseCode.InternalServerError,
		}))
		return
	}

	ctx.JSON(http.StatusOK, httpcommon.NewSuccessResponse[*[]model.SavedReceiverResponse](&savedReceivers))
}

// RenameReceiver handles the endpoint to rename a saved receiver
// @Summary Rename a saved receiver
// @Description Renames a saved receiver's nickname by receiver Id
// @Tags Receivers
// @Accept  json
// @Produce  json
// @Param id path int64 true "Receiver Id"
// @Param body body model.UpdateReceiverRequest true "Request body for renaming the receiver"
// @Success 204 "No Content"
// @Failure 400 {object} httpcommon.HttpResponse[any]
// @Failure 500 {object} httpcommon.HttpResponse[any]
// @Router /customer/saved-receiver/{id} [put]
func (handler *SavedReceiverHandler) RenameReceiver(ctx *gin.Context) {
	var updateReceiverRequest model.UpdateReceiverRequest

	if err := validation.BindJsonAndValidate(ctx, &updateReceiverRequest); err != nil {
		return
	}

	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, httpcommon.NewErrorResponse(httpcommon.Error{
			Message: "Invalid receiver Id", Field: "id", Code: httpcommon.ErrorResponseCode.MissingIdParameter,
		}))
		return
	}

	err = handler.savedReceiverService.UpdateNickname(ctx, id, updateReceiverRequest.NewNickname)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, httpcommon.NewErrorResponse(httpcommon.Error{
			Message: err.Error(), Field: "", Code: httpcommon.ErrorResponseCode.InternalServerError,
		}))
		return
	}
	ctx.AbortWithStatus(204)
}

// DeleteReceiver handles the endpoint to delete a saved receiver
// @Summary Delete a saved receiver
// @Description Deletes a saved receiver by receiver Id
// @Tags Receivers
// @Accept  json
// @Produce  json
// @Param id path int64 true "Receiver Id"
// @Success 204 "No Content"
// @Failure 400 {object} httpcommon.HttpResponse[any]
// @Failure 500 {object} httpcommon.HttpResponse[any]
// @Router /customer/saved-receiver/{id} [delete]
func (handler *SavedReceiverHandler) DeleteReceiver(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, httpcommon.NewErrorResponse(httpcommon.Error{
			Message: "Invalid receiver Id", Field: "id", Code: httpcommon.ErrorResponseCode.MissingIdParameter,
		}))
		return
	}

	err = handler.savedReceiverService.DeleteReceiver(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, httpcommon.NewErrorResponse(httpcommon.Error{
			Message: err.Error(), Field: "", Code: httpcommon.ErrorResponseCode.InternalServerError,
		}))
		return
	}
	ctx.AbortWithStatus(204)
}
