package domain

import "time"

type Employee struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Position  string    `json:"position"`
	Role      string    `json:"role"`
	Phone     string    `json:"phone"`
	Alamat    string    `json:"alamat"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
}

type PaginationMeta struct {
	Page       int `json:"page"`
	Limit      int `json:"limit"`
	Total      int `json:"total"`
	TotalPages int `json:"total_pages"`
}

type PaginatedEmployeeResponse struct {
	Data []Employee     `json:"data"`
	Meta PaginationMeta `json:"meta"`
}

type EmployeeRepository interface {
	FindAll(page, limit int) ([]Employee, int, error)
	FindByID(id int) (*Employee, error)
	Create(employee *Employee) error
	Update(employee *Employee) error
	Delete(id int) error
}

type EmployeeService interface {
	GetAllEmployees(page, limit int) (*PaginatedEmployeeResponse, error)
	GetEmployee(id int) (*Employee, error)
	CreateEmployee(employee *Employee) error
	UpdateEmployee(employee *Employee) error
	DeleteEmployee(id int) error
}
