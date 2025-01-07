package v1

import (
	"net/http"
	"strconv"

	"github.com/21CLC01-WNC-Banking/WNC-Banking-BE/internal/controller/http/middleware"
	"github.com/21CLC01-WNC-Banking/WNC-Banking-BE/internal/domain/entity"
	httpcommon "github.com/21CLC01-WNC-Banking/WNC-Banking-BE/internal/domain/http_common"
	"github.com/21CLC01-WNC-Banking/WNC-Banking-BE/internal/service"
	"github.com/gin-gonic/gin"
)

type CustomerHandler struct {
	notificationService service.NotificationService
	transactionService  service.TransactionService
}

func NewCustomerHandler(notificationService service.NotificationService, transactionService service.TransactionService) *CustomerHandler {
	return &CustomerHandler{
		notificationService: notificationService,
		transactionService:  transactionService,
	}
}

// @Summary Seen notification by id
// @Description Seen notification by id
// @Tags Customer
// @Param id query int64 true "Id of notification"
// @Produce  json
// @Router /customer/notification [patch]
// @Param  Authorization header string true "Authorization: Bearer"
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
// @Param  Authorization header string true "Authorization: Bearer"
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
		return
	}
	if len(notifications) == 0 {
		notifications = make([]entity.Notification, 0)
	}
	c.JSON(http.StatusOK, httpcommon.NewSuccessResponse(&notifications))
}

// @Summary Get All Transactions
// @Description Get All Transactions
// @Tags Customer
// @Produce  json
// @Router /customer/transaction [get]
// @Param  Authorization header string true "Authorization: Bearer"
// @Success 200 {object} httpcommon.HttpResponse[[]entity.Transaction]
// @Failure 400 {object} httpcommon.HttpResponse[any]
// @Failure 500 {object} httpcommon.HttpResponse[any]
func (h *CustomerHandler) GetTransactions(c *gin.Context) {
	customerId := middleware.GetUserIdHelper(c)

	transactions, err := h.transactionService.GetTransactionsByCustomerId(c, customerId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, httpcommon.NewErrorResponse(
			httpcommon.Error{
				Message: err.Error(),
			},
		))
		return
	}
	c.JSON(http.StatusOK, httpcommon.NewSuccessResponse(&transactions))
}

// @Summary Get Transaction By Id
// @Description Retrieve a specific transaction by its Id for a given customer.
// @Tags Customer
// @Produce json
// @Param transactionId path string true "Transaction Id"
// @Router /customer/transaction/{transactionId} [get]
// @Param  Authorization header string true "Authorization: Bearer"
// @Success 200 {object} httpcommon.HttpResponse[entity.Transaction]
// @Failure 400 {object} httpcommon.HttpResponse[any]
// @Failure 500 {object} httpcommon.HttpResponse[any]
func (h *CustomerHandler) GetTransactionById(c *gin.Context) {
	customerId := middleware.GetUserIdHelper(c)
	transactionId := c.Param("transactionId")

	transaction, err := h.transactionService.GetTransactionByIdAndCustomerId(c, customerId, transactionId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, httpcommon.NewErrorResponse(
			httpcommon.Error{
				Message: err.Error(),
			},
		))
		return
	}
	c.JSON(http.StatusOK, httpcommon.NewSuccessResponse(&transaction))
}
