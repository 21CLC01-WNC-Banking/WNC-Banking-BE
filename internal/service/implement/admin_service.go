package serviceimplement

import (
	"github.com/21CLC01-WNC-Banking/WNC-Banking-BE/internal/bean"
	"github.com/21CLC01-WNC-Banking/WNC-Banking-BE/internal/domain/entity"
	"github.com/21CLC01-WNC-Banking/WNC-Banking-BE/internal/domain/model"
	"github.com/21CLC01-WNC-Banking/WNC-Banking-BE/internal/repository"
	"github.com/21CLC01-WNC-Banking/WNC-Banking-BE/internal/service"
	"github.com/gin-gonic/gin"
)

type AdminService struct {
	staffRepository       repository.StaffRepository
	passwordEncoder       bean.PasswordEncoder
	transactionRepository repository.TransactionRepository
}

func NewAdminService(staffRepository repository.StaffRepository, passwordEncoder bean.PasswordEncoder, transactionRepository repository.TransactionRepository) service.AdminService {
	return &AdminService{staffRepository: staffRepository, passwordEncoder: passwordEncoder, transactionRepository: transactionRepository}
}

func (a *AdminService) GetAllStaff(ctx *gin.Context) ([]entity.User, error) {
	return a.staffRepository.GetAll(ctx)
}

func (a *AdminService) GetOneStaff(ctx *gin.Context, staffId int64) (*entity.User, error) {
	return a.staffRepository.GetOneById(ctx, staffId)
}

func (a *AdminService) CreateOneStaff(ctx *gin.Context, request *model.CreateStaffRequest) (int64, error) {
	hashedPassword, err := a.passwordEncoder.Encrypt(request.Password)
	if err != nil {
		return 0, err
	}

	return a.staffRepository.CreateOne(ctx, &entity.User{
		Email:       request.Email,
		Name:        request.Name,
		PhoneNumber: request.PhoneNumber,
		Password:    hashedPassword,
	})
}

func (a *AdminService) DeleteOneStaff(ctx *gin.Context, staffId int64) error {
	return a.staffRepository.DeleteOne(ctx, staffId)
}

func (a *AdminService) UpdateOneStaff(ctx *gin.Context, request *model.UpdateStaffRequest) error {
	if request.Password != "" {
		hashedPassword, err := a.passwordEncoder.Encrypt(request.Password)
		if err != nil {
			return err
		}
		request.Password = hashedPassword
	}

	return a.staffRepository.UpdateOneStaff(ctx, &entity.User{
		Id:          request.Id,
		Email:       request.Email,
		Name:        request.Name,
		Password:    request.Password,
		PhoneNumber: request.PhoneNumber,
	})
}

func (a *AdminService) GetExternalTransactions(ctx *gin.Context, filter model.GetExternalTransactionRequest) ([]entity.Transaction, error) {
	return a.transactionRepository.GetExternalTransactionsWithFilter(ctx, filter.FromDate, filter.ToDate, filter.BankId)
}
