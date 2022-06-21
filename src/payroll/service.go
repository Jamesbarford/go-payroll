/*
 * These are presently
 */
package payroll

import (
	"log"
)

type PayrollService struct {
	logger     *log.Logger
	repository *PayrollRepository
}

func NewPayrollService(repository *PayrollRepository) *PayrollService {
	return &PayrollService{
		logger:     log.Default(),
		repository: repository,
	}
}

func (service *PayrollService) GetPayrollForUser(userId int, month int) (string, error) {
	return service.repository.GetPayrollForUser(userId, month)
}

func (service *PayrollService) AddUserToPayRoll(request *PayrollAddUserRequest) error {
	return service.repository.AddUserToPayRoll(request)
}
