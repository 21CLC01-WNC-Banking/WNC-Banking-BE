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

type AdminHandler struct {
	adminService       service.AdminService
	partnerBankService service.PartnerBankService
}

func NewAdminHandler(adminService service.AdminService, partnerBankService service.PartnerBankService) *AdminHandler {
	return &AdminHandler{
		adminService:       adminService,
		partnerBankService: partnerBankService,
	}
}

// @Summary Admin get all staff
// @Description Admin get all staff
// @Tags Admins
// @Produce  json
// @Router /admin/staff [get]
// @Param  Authorization header string true "Authorization: Bearer"
// @Success 200 {object} httpcommon.HttpResponse[[]entity.User]
// @Failure 500 {object} httpcommon.HttpResponse[any]
func (handler *AdminHandler) GetAllStaff(c *gin.Context) {
	staffs, err := handler.adminService.GetAllStaff(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, httpcommon.NewErrorResponse(
			httpcommon.Error{Message: err.Error()},
		))
		return
	}
	c.JSON(http.StatusOK, httpcommon.NewSuccessResponse(&staffs))
}

// @Summary Admin get one staff
// @Description Admin get one staff
// @Tags Admins
// @Produce  json
// @Param staffId path int64 true "Staff Id"
// @Router /admin/staff/{staffId} [get]
// @Param  Authorization header string true "Authorization: Bearer"
// @Success 200 {object} httpcommon.HttpResponse[entity.User]
// @Failure 400 {object} httpcommon.HttpResponse[any]
// @Failure 500 {object} httpcommon.HttpResponse[any]
func (handler *AdminHandler) GetOneStaff(c *gin.Context) {
	staffIdStr := c.Param("staffId")
	staffId, err := strconv.ParseInt(staffIdStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, httpcommon.NewErrorResponse(
			httpcommon.Error{
				Message: err.Error(),
				Code:    httpcommon.ErrorResponseCode.InvalidRequest,
			},
		))
		return
	}

	staff, err := handler.adminService.GetOneStaff(c, staffId)
	if err != nil {
		if err.Error() == httpcommon.ErrorMessage.SqlxNoRow {
			c.JSON(http.StatusNotFound, httpcommon.NewErrorResponse(
				httpcommon.Error{Message: err.Error()},
			))
			return
		}
		c.JSON(http.StatusInternalServerError, httpcommon.NewErrorResponse(
			httpcommon.Error{Message: err.Error()},
		))
		return
	}
	c.JSON(http.StatusOK, httpcommon.NewSuccessResponse(&staff))
}

// @Summary Admin create one staff
// @Description Admin create one staff
// @Tags Admins
// @Produce  json
// @Param request body model.CreateStaffRequest true "Staff payload"
// @Router /admin/staff [post]
// @Param  Authorization header string true "Authorization: Bearer"
// @Success 200 {object} httpcommon.HttpResponse[model.CreateStaffResponse]
// @Failure 400 {object} httpcommon.HttpResponse[any]
// @Failure 500 {object} httpcommon.HttpResponse[any]
func (handler *AdminHandler) CreateOneStaff(c *gin.Context) {
	var request model.CreateStaffRequest

	err := validation.BindJsonAndValidate(c, &request)
	if err != nil {
		return
	}

	id, err := handler.adminService.CreateOneStaff(c, &request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, httpcommon.NewErrorResponse(
			httpcommon.Error{Message: err.Error()},
		))
		return
	}

	c.JSON(http.StatusOK, httpcommon.NewSuccessResponse(&model.CreateStaffResponse{Id: id}))
}

// @Summary Admin delete one staff
// @Description Admin delete one staff
// @Tags Admins
// @Produce  json
// @Param staffId path int64 true "Staff Id"
// @Router /admin/staff/{staffId} [delete]
// @Param  Authorization header string true "Authorization: Bearer"
// @Success 204 "No content"
// @Failure 400 {object} httpcommon.HttpResponse[any]
// @Failure 500 {object} httpcommon.HttpResponse[any]
func (handler *AdminHandler) DeleteOneStaff(c *gin.Context) {
	staffIdStr := c.Param("staffId")
	staffId, err := strconv.ParseInt(staffIdStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, httpcommon.NewErrorResponse(
			httpcommon.Error{
				Message: err.Error(),
				Code:    httpcommon.ErrorResponseCode.InvalidRequest,
			},
		))
		return
	}

	err = handler.adminService.DeleteOneStaff(c, staffId)
	if err != nil {
		if err.Error() == httpcommon.ErrorMessage.SqlxNoRow {
			c.JSON(http.StatusNotFound, httpcommon.NewErrorResponse(
				httpcommon.Error{Message: err.Error()},
			))
			return
		}
		c.JSON(http.StatusInternalServerError, httpcommon.NewErrorResponse(
			httpcommon.Error{Message: err.Error()},
		))
		return
	}
	c.AbortWithStatus(http.StatusNoContent)
}

// @Summary Admin update one staff ( only update non-empty field )
// @Description Admin update one staff
// @Tags Admins
// @Produce  json
// @Param request body model.UpdateStaffRequest true "Staff payload"
// @Router /admin/staff [put]
// @Param  Authorization header string true "Authorization: Bearer"
// @Success 204 "No content"
// @Failure 400 {object} httpcommon.HttpResponse[any]
// @Failure 500 {object} httpcommon.HttpResponse[any]
func (handler *AdminHandler) UpdateOneStaff(c *gin.Context) {
	var request model.UpdateStaffRequest

	err := validation.BindJsonAndValidate(c, &request)
	if err != nil {
		return
	}

	err = handler.adminService.UpdateOneStaff(c, &request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, httpcommon.NewErrorResponse(
			httpcommon.Error{Message: err.Error()},
		))
		return
	}

	c.AbortWithStatus(http.StatusNoContent)
}

// @Summary Add Partner Bank
// @Description Add a partner bank
// @Tags Admins
// @Accept json
// @Param request body model.PartnerBankRequest true "PartnerBank payload"
// @Produce  json
// @Router /admin/partner-bank [post]
// @Param  Authorization header string true "Authorization: Bearer"
// @Success 200 "No Content"
// @Failure 400 {object} httpcommon.HttpResponse[any]
// @Failure 500 {object} httpcommon.HttpResponse[any]
func (handler *AdminHandler) AddPartnerBank(c *gin.Context) {
	var partnerRequest model.PartnerBankRequest
	if err := validation.BindJsonAndValidate(c, &partnerRequest); err != nil {
		return
	}
	err := handler.partnerBankService.AddPartnerBank(c, partnerRequest)
	if err != nil {
		c.JSON(http.StatusInternalServerError, httpcommon.NewErrorResponse(httpcommon.Error{
			Message: err.Error(), Field: "", Code: httpcommon.ErrorResponseCode.InternalServerError,
		}))
		return
	}
	c.AbortWithStatus(200)
}

// @Summary Get external transactions
// @Description Get external transactions
// @Tags Admins
// @Accept json
// @Param fromDate query string true "yyyy-MM-dd"
// @Param toDate query string true "yyyy-MM-dd"
// @Param bankId query int64 false "Partner Bank Id"
// @Produce  json
// @Router /admin/external-transaction [get]
// @Param  Authorization header string true "Authorization: Bearer"
// @Success 200 {object} httpcommon.HttpResponse[entity.Transaction]
// @Failure 400 {object} httpcommon.HttpResponse[any]
// @Failure 500 {object} httpcommon.HttpResponse[any]
func (handler *AdminHandler) GetExternalTransactions(c *gin.Context) {
	BankIdStr := c.Param("BankId")
	bankId, _ := strconv.ParseInt(BankIdStr, 10, 64)
	req := model.GetExternalTransactionRequest{
		BankId:   bankId,
		FromDate: c.Query("fromDate"),
		ToDate:   c.Query("toDate"),
	}

	transactions, err := handler.adminService.GetExternalTransactions(c, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, httpcommon.NewErrorResponse(
			httpcommon.Error{Message: err.Error()},
		))
		return
	}

	if len(transactions) == 0 {
		transactions = make([]model.GetExternalTransactionResponse, 0)
	}
	c.JSON(http.StatusOK, httpcommon.NewSuccessResponse(&transactions))
}
