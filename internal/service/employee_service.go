package service

import (
	"errors"
	"regexp"
	"strings"

	"karyawan-app/internal/domain"
)

var (
	emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
)

type employeeService struct {
	repo domain.EmployeeRepository
}

func NewEmployeeService(repo domain.EmployeeRepository) domain.EmployeeService {
	return &employeeService{repo: repo}
}

func (s *employeeService) GetAllEmployees() ([]domain.Employee, error) {
	return s.repo.FindAll()
}

func (s *employeeService) GetEmployee(id int) (*domain.Employee, error) {
	return s.repo.FindByID(id)
}

func (s *employeeService) CreateEmployee(employee *domain.Employee) error {
	if err := validateEmployee(employee); err != nil {
		return err
	}
	return s.repo.Create(employee)
}

func (s *employeeService) UpdateEmployee(employee *domain.Employee) error {
	if employee.ID == 0 {
		return errors.New("employee ID is required")
	}
	if err := validateEmployee(employee); err != nil {
		return err
	}
	return s.repo.Update(employee)
}

func (s *employeeService) DeleteEmployee(id int) error {
	return s.repo.Delete(id)
}

func validateEmployee(employee *domain.Employee) error {
	if strings.TrimSpace(employee.Name) == "" {
		return errors.New("name is required")
	}

	if !isValidEmail(employee.Email) {
		return errors.New("invalid email format")
	}

	if strings.TrimSpace(employee.Role) == "" {
		return errors.New("role is required")
	}

	if strings.TrimSpace(employee.Phone) == "" {
		return errors.New("phone is required")
	}

	if strings.TrimSpace(employee.Alamat) == "" {
		return errors.New("alamat is required")
	}

	return nil
}

func isValidEmail(email string) bool {
	return emailRegex.MatchString(email)
}
