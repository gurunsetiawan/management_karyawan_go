package service

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
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

func (s *employeeService) GetAllEmployees(page, limit int, search, status string) (*domain.PaginatedEmployeeResponse, error) {
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 10
	}

	data, total, err := s.repo.FindAll(page, limit, search, status)
	if err != nil {
		return nil, err
	}

	totalPages := total / limit
	if total%limit > 0 {
		totalPages++
	}

	return &domain.PaginatedEmployeeResponse{
		Data: data,
		Meta: domain.PaginationMeta{
			Page:       page,
			Limit:      limit,
			Total:      total,
			TotalPages: totalPages,
		},
	}, nil
}

func (s *employeeService) GetEmployee(id int) (*domain.Employee, error) {
	return s.repo.FindByID(id)
}

func (s *employeeService) CreateEmployee(employee *domain.Employee) error {
	// Sanitize input
	sanitizeEmployee(employee)
	
	if err := validateEmployee(employee); err != nil {
		return err
	}
	return s.repo.Create(employee)
}

func (s *employeeService) UpdateEmployee(employee *domain.Employee) error {
	if employee.ID == 0 {
		return errors.New("employee ID is required")
	}
	
	// Sanitize input
	sanitizeEmployee(employee)
	
	if err := validateEmployee(employee); err != nil {
		return err
	}
	return s.repo.Update(employee)
}

func (s *employeeService) DeleteEmployee(id int) error {
	return s.repo.Delete(id)
}

// sanitizeEmployee sanitizes all string fields in the employee struct
func sanitizeEmployee(employee *domain.Employee) {
	employee.Name = sanitizeInput(employee.Name)
	employee.Email = sanitizeInput(employee.Email)
	employee.Position = sanitizeInput(employee.Position)
	employee.Role = sanitizeInput(employee.Role)
	employee.Phone = sanitizeInput(employee.Phone)
	employee.Alamat = sanitizeInput(employee.Alamat)
}

// sanitizeInput removes HTML tags and dangerous characters to prevent XSS
func sanitizeInput(input string) string {
	input = strings.TrimSpace(input)
	input = strings.ReplaceAll(input, "<", "&lt;")
	input = strings.ReplaceAll(input, ">", "&gt;")
	input = strings.ReplaceAll(input, "\"", "&quot;")
	input = strings.ReplaceAll(input, "'", "&#39;")
	return input
}

func validateEmployee(employee *domain.Employee) error {
	if strings.TrimSpace(employee.Name) == "" {
		return errors.New("name is required")
	}

	if !isValidEmail(employee.Email) {
		return errors.New("invalid email format")
	}

	if strings.TrimSpace(employee.Position) == "" {
		return errors.New("position is required")
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

// ImportCSV parses a CSV file and imports valid rows.
// Returns success count, a list of failures, and error if the file itself is invalid.
func (s *employeeService) ImportCSV(csvData io.Reader) (int, []string, error) {
	reader := csv.NewReader(csvData)
	
	// Read header
	header, err := reader.Read()
	if err != nil {
		if err == io.EOF {
			return 0, nil, errors.New("file CSV kosong")
		}
		return 0, nil, errors.New("gagal membaca header CSV: " + err.Error())
	}
	
	// Basic validation of header
	if len(header) < 6 {
		return 0, nil, errors.New("format CSV tidak valid: kurang kolom")
	}

	var successCount int
	var failures []string
	rowNum := 1 // Header is row 1

	for {
		row, err := reader.Read()
		if err == io.EOF {
			break
		}
		rowNum++
		
		if err != nil {
			failures = append(failures, fmt.Sprintf("Baris %d: Gagal dibaca - %v", rowNum, err))
			continue
		}

		if len(row) < 6 {
			failures = append(failures, fmt.Sprintf("Baris %d: Jumlah kolom tidak sesuai", rowNum))
			continue
		}

		emp := &domain.Employee{
			Name:     strings.TrimSpace(row[0]),
			Email:    strings.TrimSpace(row[1]),
			Position: strings.TrimSpace(row[2]),
			Role:     strings.TrimSpace(row[3]),
			Phone:    strings.TrimSpace(row[4]),
			Alamat:   strings.TrimSpace(row[5]),
		}

		if err := validateEmployee(emp); err != nil {
			failures = append(failures, fmt.Sprintf("Baris %d: %v", rowNum, err))
			continue
		}

		if err := s.repo.Create(emp); err != nil {
			failures = append(failures, fmt.Sprintf("Baris %d: Gagal disimpan - %v", rowNum, err))
			continue
		}

		successCount++
	}

	return successCount, failures, nil
}

func (s *employeeService) ExportCSV() ([]byte, error) {
	return s.repo.ExportCSV()
}
