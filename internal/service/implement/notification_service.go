package serviceimplement

import (
	"errors"
	"github.com/21CLC01-WNC-Banking/WNC-Banking-BE/internal/domain/entity"
	httpcommon "github.com/21CLC01-WNC-Banking/WNC-Banking-BE/internal/domain/http_common"
	"github.com/21CLC01-WNC-Banking/WNC-Banking-BE/internal/repository"
	"github.com/21CLC01-WNC-Banking/WNC-Banking-BE/internal/service"
	"github.com/gin-gonic/gin"
)

type NotificationService struct {
	notificationRepository repository.NotificationRepository
}

func NewNotificationService(
	notificationRepository repository.NotificationRepository,
) service.NotificationService {
	return &NotificationService{
		notificationRepository: notificationRepository,
	}
}

func (n NotificationService) GetAllNotifications(ctx *gin.Context, userId int64) ([]entity.Notification, error) {
	return n.notificationRepository.GetAllByUserId(ctx, userId)
}

func (n NotificationService) PatchNotification(ctx *gin.Context, userId int64, id int64) error {
	userIdFromNotification := n.notificationRepository.GetCustomerIdById(ctx, id)
	if userIdFromNotification == 0 || userIdFromNotification != userId {
		return errors.New(httpcommon.ErrorMessage.BadCredential)
	}

	return n.notificationRepository.PatchSeenById(ctx, id)
}
