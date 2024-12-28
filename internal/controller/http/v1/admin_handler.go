package v1

import (
	httpcommon "github.com/21CLC01-WNC-Banking/WNC-Banking-BE/internal/domain/http_common"
	"github.com/21CLC01-WNC-Banking/WNC-Banking-BE/internal/service"
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
// @Param staffId path int64 true "Staff ID"
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
