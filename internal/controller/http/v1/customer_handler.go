package v1

import (
	"github.com/21CLC01-WNC-Banking/WNC-Banking-BE/internal/controller/http/middleware"
	"github.com/21CLC01-WNC-Banking/WNC-Banking-BE/internal/domain/entity"
	httpcommon "github.com/21CLC01-WNC-Banking/WNC-Banking-BE/internal/domain/http_common"
	"github.com/21CLC01-WNC-Banking/WNC-Banking-BE/internal/service"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type CustomerHandler struct {
	notificationService service.NotificationService
}

func NewCustomerHandler(notificationService service.NotificationService) *CustomerHandler {
	return &CustomerHandler{
		notificationService: notificationService,
	}
}

// @Summary Seen notification by id
// @Description Seen notification by id
// @Tags Customer
// @Param id query int64 true "Id of notification"
// @Produce  json
// @Router /customer/notification [patch]
// @Success 204 "No content"
// @Failure 400 {object} httpcommon.HttpResponse[any]
// @Failure 500 {object} httpcommon.HttpResponse[any]
func (h *CustomerHandler) SeenNotification(c *gin.Context) {
	customerId := middleware.GetUserIdHelper(c)
	notificationIdStr := c.Param(`notificationId`)

	notificationId, err := strconv.ParseInt(notificationIdStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, httpcommon.NewErrorResponse(
			httpcommon.Error{Message: err.Error()},
		))
		return
	}

	err = h.notificationService.PatchNotification(c, customerId, notificationId)
	if err != nil {
		c.JSON(http.StatusBadRequest, httpcommon.NewErrorResponse(
			httpcommon.Error{Message: err.Error()},
		))
		return
	}
	c.AbortWithStatus(http.StatusNoContent)
}

// @Summary Get All Notification
// @Description Get All Notification
// @Tags Customer
// @Produce  json
// @Router /customer/notification [get]
// @Success 200 {object} httpcommon.HttpResponse[[]entity.Notification]
// @Failure 400 {object} httpcommon.HttpResponse[any]
// @Failure 500 {object} httpcommon.HttpResponse[any]
func (h *CustomerHandler) GetNotifications(c *gin.Context) {
	customerId := middleware.GetUserIdHelper(c)

	notifications, err := h.notificationService.GetAllNotifications(c, customerId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, httpcommon.NewErrorResponse(
			httpcommon.Error{
				Message: err.Error(),
			},
		))
	}
	if len(notifications) == 0 {
		notifications = make([]entity.Notification, 0)
	}
	c.JSON(http.StatusOK, httpcommon.NewSuccessResponse(&notifications))
}