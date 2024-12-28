package v1

import (
	httpcommon "github.com/21CLC01-WNC-Banking/WNC-Banking-BE/internal/domain/http_common"
	"github.com/21CLC01-WNC-Banking/WNC-Banking-BE/internal/domain/model"
	"github.com/21CLC01-WNC-Banking/WNC-Banking-BE/internal/service"
	"github.com/21CLC01-WNC-Banking/WNC-Banking-BE/internal/utils/validation"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type AdminHandler struct {
	adminService service.AdminService
}

func NewAdminHandler(adminService service.AdminService) *AdminHandler {
	return &AdminHandler{
		adminService: adminService,
	}
}

// @Summary Admin get all staff
// @Description Admin get all staff
// @Tags Admins
// @Produce  json
// @Router /admin/staff [get]
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
