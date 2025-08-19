package domain

import "time"

type Employee struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Role      string    `json:"role"`
	Phone     string    `json:"phone"`
	Alamat    string    `json:"alamat"`
	CreatedAt time.Time `json:"created_at"`
}

type EmployeeRepository interface {
	FindAll() ([]Employee, error)
	FindByID(id int) (*Employee, error)
	Create(employee *Employee) error
	Update(employee *Employee) error
	Delete(id int) error
}

type EmployeeService interface {
	GetAllEmployees() ([]Employee, error)
	GetEmployee(id int) (*Employee, error)
	CreateEmployee(employee *Employee) error
	UpdateEmployee(employee *Employee) error
	DeleteEmployee(id int) error
}
